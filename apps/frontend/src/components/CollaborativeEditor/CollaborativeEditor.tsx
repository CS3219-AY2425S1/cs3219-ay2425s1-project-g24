// Referenced from example in https://www.npmjs.com/package/y-codemirror.next
import React, {
  Dispatch,
  ForwardedRef,
  forwardRef,
  RefObject,
  SetStateAction,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
  MutableRefObject,
} from "react";
import * as Y from "yjs";
import { yCollab } from "y-codemirror.next";
import { WebrtcProvider } from "y-webrtc";
import { EditorView, basicSetup } from "codemirror";
import { keymap } from "@codemirror/view";
import { indentWithTab } from "@codemirror/commands";
import { EditorState, Compartment } from "@codemirror/state";
import { javascript, javascriptLanguage } from "@codemirror/lang-javascript";
import { python, pythonLanguage } from "@codemirror/lang-python";
import { javaLanguage, java } from "@codemirror/lang-java";
import { cppLanguage, cpp } from "@codemirror/lang-cpp";
import { goLanguage, go } from "@codemirror/lang-go";
import "./styles.scss";
import { message, Select } from "antd";
import { language } from "@codemirror/language";
import { ProgrammingLanguageOptions } from "@/utils/SelectOptions";
import { ExecutionResults, SubmissionResults } from "@/app/services/execute";

interface CollaborativeEditorProps {
  user: string;
  collaborationId: string;
  updateLanguage: (language: string) => void;
  setMatchedUser: Dispatch<SetStateAction<string>>;
  handleCloseCollaboration: (type: string) => void;
  providerRef: MutableRefObject<WebrtcProvider | null>;
  matchedUser: string;
  onCodeChange: (code: string) => void;
  updateSubmissionResults: (results: SubmissionResults) => void;
  updateExecutionResults: (results: ExecutionResults) => void;
  updateExecuting: (executing: boolean) => void;
  updateSubmitting: (submitting: boolean) => void;
}

export interface CollaborativeEditorHandle {
  endSession: () => void;
}

interface AwarenessUpdate {
  added: number[];
  updated: number[];
  removed: number[];
}

interface Awareness {
  sessionEnded: boolean;
  user: {
    name: string;
    color: string;
    colorLight: string;
  };
  submissionResultsState: {
    submissionResults: SubmissionResults;
    id: number;
  };
  executionResultsState: {
    executionResults: ExecutionResults;
    id: number;
  };
  executingState: {
    executing: boolean;
    id: number;
  };
  submittingState: {
    submitting: boolean;
    id: number;
  };
  editorState: {
    language: string;
    id: number;
  };
}

export const usercolors = [
  { color: "#30bced", light: "#30bced33" },
  { color: "#6eeb83", light: "#6eeb8333" },
  { color: "#ffbc42", light: "#ffbc4233" },
  { color: "#ecd444", light: "#ecd44433" },
  { color: "#ee6352", light: "#ee635233" },
  { color: "#9ac2c9", light: "#9ac2c933" },
  { color: "#8acb88", light: "#8acb8833" },
  { color: "#1be7ff", light: "#1be7ff33" },
];

// Retrieve random colour
export const userColor =
  usercolors[Math.floor(Math.random() * (usercolors.length - 1))];

const CollaborativeEditor = forwardRef(
  (
    props: CollaborativeEditorProps,
    ref: ForwardedRef<CollaborativeEditorHandle>
  ) => {
    const editorRef = useRef(null);
    // const providerRef = useRef<WebrtcProvider | null>(null);
    const [selectedLanguage, setSelectedLanguage] = useState("python");
    const [mounted, setMounted] = useState(false);
    let sessionEndNotified = false;
    let sessionEndTimeout: any;

    const languageConf = new Compartment();

    const codeChangeListener = EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        props.onCodeChange(update.state.doc.toString());
      }
    });

    // Referenced: https://codemirror.net/examples/config/#dynamic-configuration
    const autoLanguage = EditorState.transactionExtender.of((tr) => {
      if (!tr.docChanged) return null;
      const editorLanguage = localStorage.getItem("editor-language") ?? "";
      let stateIsJs = tr.startState.facet(language) == javascriptLanguage;
      let stateIsPython = tr.startState.facet(language) == pythonLanguage;
      if (
        (stateIsJs && editorLanguage.toLowerCase() === "javascript") ||
        (stateIsPython && editorLanguage.toLowerCase() === "python")
      )
        return null;

      return {
        effects: languageConf.reconfigure(
          editorLanguage.toLowerCase() === "javascript"
            ? javascript()
            : python()
        ),
      };
    });

    const [messageApi, contextHolder] = message.useMessage();

    const success = (message: string) => {
      messageApi.open({
        type: "success",
        content: message,
      });
    };

    const info = (message: string) => {
      messageApi.open({
        type: "info",
        content: message,
      });
    };

    const error = (message: string) => {
      messageApi.open({
        type: "error",
        content: message,
      });
    };

    const warning = (message: string) => {
      messageApi.open({
        type: "warning",
        content: message,
      });
    };

    let latestExecutionId: number = new Date(0).getTime();
    let latestSubmissionId: number = new Date(0).getTime();
    let latestExecutingId: number = new Date(0).getTime();
    let latestSubmittingId: number = new Date(0).getTime();
    let latestLanguageChangeId: number = new Date(0).getTime();

    useImperativeHandle(ref, () => ({
      endSession: () => {
        if (props.providerRef.current) {
          // Set awareness state to indicate session ended to notify peer about session ending
          props.providerRef.current.awareness.setLocalStateField(
            "sessionEnded",
            true
          );
          success("Session ended. All participants will be notified.");

          // Event when peer disconnected, in between 0-15s user decides to end session before auto end session is triggered from inactivity
          if (sessionEndTimeout) {
            clearTimeout(sessionEndTimeout);
            sessionEndTimeout = null;
          }
        }
      },
    }));

    useEffect(() => {
      localStorage.setItem("editor-language", selectedLanguage);

      if (props.providerRef.current && mounted) {
        latestLanguageChangeId = Date.now();
        props.providerRef.current.awareness.setLocalStateField("editorState", {
          language: selectedLanguage,
          id: latestLanguageChangeId,
        });
        props.updateLanguage(selectedLanguage);
        success(`Changed Code Editor's language to ${selectedLanguage}`);
      } else {
        setMounted(true);
      }
    }, [selectedLanguage]);

    useEffect(() => {
      if (process.env.NEXT_PUBLIC_SIGNALLING_SERVICE_URL === undefined) {
        error("Missing Signalling Service Url");
        return;
      }

      const ydoc = new Y.Doc();
      const provider = new WebrtcProvider(props.collaborationId, ydoc, {
        signaling: [process.env.NEXT_PUBLIC_SIGNALLING_SERVICE_URL],
        maxConns: 2,
      });

      props.providerRef.current = provider;
      const ytext = ydoc.getText("codemirror");
      const undoManager = new Y.UndoManager(ytext);

      provider.awareness.setLocalStateField("user", {
        name: props.user,
        color: userColor.color,
        colorLight: userColor.light,
      });

      // Check initial awareness states
      const states = provider.awareness.getStates();
      for (const [clientID, state] of Array.from(states)) {
        if (state.user && state.user.name !== props.user) {
          props.setMatchedUser(state.user.name);
          break;
        }
      }

      // Listen for awareness changes
      provider.awareness.on("change", () => {
        const updatedStates = provider.awareness.getStates();

        // Check the length of updatedStates, if it is equal to 1, we trigger endsession in 15s but if updated to 2, cancel endsession
        if (sessionEndTimeout && updatedStates.size == 2) {
          clearTimeout(sessionEndTimeout);
          sessionEndTimeout = null;
        }

        // If there's only one participant, set a timeout for 15 seconds
        if (updatedStates.size === 1) {
          sessionEndTimeout = setTimeout(() => {
            // Trigger end session logic here
            info(
              `Session has ended due to inactivity from matched user ${props.matchedUser}`
            );
            props.handleCloseCollaboration("peer");
            if (props.providerRef.current) {
              props.providerRef.current.disconnect();
            }
          }, 15000); // 15 seconds
        }

        for (const [clientID, state] of Array.from(updatedStates)) {
          if (state.sessionEnded && state.user.name !== props.user) {
            if (!sessionEndNotified) {
              info(
                `Session has been ended by another participant ${state.user.name}`
              );

              props.handleCloseCollaboration("peer");
              sessionEndNotified = true;
              if (props.providerRef.current) {
                props.providerRef.current.disconnect();
              }
              return;
            }
          }

          if (state.user && state.user.name !== props.user) {
            props.setMatchedUser(state.user.name);
            break;
          }
        }
      });

      // Listener for awareness updates to receive submission results from peer
      provider.awareness.on("update", ({ added, updated }: AwarenessUpdate) => {
        added
          .concat(updated)
          .filter((clientId) => clientId !== provider.awareness.clientID)
          .forEach((clientID) => {
            const state = provider.awareness
              .getStates()
              .get(clientID) as Awareness;

            // New section to check for changes in language
            if (
              state &&
              state.editorState &&
              state.editorState.id !== latestLanguageChangeId
            ) {
              latestSubmissionId = state.editorState.id;
              setSelectedLanguage(state.editorState.language);
              // if (props.user === state.user.name) {
              //   console.log("ownself update ownself");
              // } else {
              //   console.log("others update ownself");
              // }
            }

            if (
              state &&
              state.submissionResultsState &&
              state.submissionResultsState.id !== latestSubmissionId
            ) {
              latestSubmissionId = state.submissionResultsState.id;
              props.updateSubmissionResults(
                state.submissionResultsState.submissionResults
              );
              messageApi.open({
                type: "success",
                content: `${
                  props.matchedUser ?? "Peer"
                } saved code successfully!`,
              });
            }

            if (
              state &&
              state.executionResultsState &&
              state.executionResultsState.id !== latestExecutionId
            ) {
              latestExecutionId = state.executionResultsState.id;
              props.updateExecutionResults(
                state.executionResultsState.executionResults
              );
              messageApi.open({
                type: "success",
                content: `${
                  props.matchedUser ?? "Peer"
                } ran test cases. Review the results below.`,
              });
            }

            if (
              state &&
              state.executingState &&
              state.executingState.id !== latestExecutingId
            ) {
              latestExecutingId = state.executingState.id;
              props.updateExecuting(state.executingState.executing);
              if (state.executingState.executing) {
                messageApi.open({
                  type: "info",
                  content: `${
                    props.matchedUser ?? "Peer"
                  } is running test cases...`,
                });
              }
            }

            if (
              state &&
              state.submittingState &&
              state.submittingState.id !== latestSubmittingId
            ) {
              latestSubmittingId = state.submittingState.id;
              props.updateSubmitting(state.submittingState.submitting);
              if (state.submittingState.submitting) {
                messageApi.open({
                  type: "info",
                  content: `${props.matchedUser ?? "Peer"} is saving code...`,
                });
              }
            }
          });
      });

      const state = EditorState.create({
        doc: ytext.toString(),
        extensions: [
          basicSetup,
          languageConf.of(python()),
          // languageConf.of(node()),
          autoLanguage,
          yCollab(ytext, provider.awareness, { undoManager }),
          keymap.of([indentWithTab]),
          codeChangeListener,
          EditorView.theme({
            "&": { height: "100%", overflow: "hidden" }, // Enable Scroll
          }),
        ],
      });

      const view = new EditorView({
        state,
        parent: editorRef.current || undefined,
      });

      return () => {
        // Cleanup on component unmount
        view.destroy();
        provider.disconnect();
        ydoc.destroy();
      };
    }, []);

    return (
      <>
        {contextHolder}
        <div className="code-second-container">
          <div className="code-language">Select Language:</div>
          <Select
            className="language-select"
            // defaultValue={localStorage.getItem("editor-language")}
            value={selectedLanguage}
            options={ProgrammingLanguageOptions}
            onSelect={(val) => {
              setSelectedLanguage(val);
            }}
          />
        </div>
        <div
          ref={editorRef}
          style={{ height: "400px", border: "1px solid #ddd" }}
        />
        {/* <div className="language-detected">
          <strong>Current Language Detected: </strong> {selectedLanguage}
        </div> */}
      </>
    );
  }
);

export default CollaborativeEditor;

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
import { keymap } from "@codemirror/view"
import { indentWithTab } from "@codemirror/commands"
import { EditorState, Compartment, Extension } from "@codemirror/state";
import { javascript, javascriptLanguage } from "@codemirror/lang-javascript";
import { python, pythonLanguage } from "@codemirror/lang-python";
import { javaLanguage, java } from "@codemirror/lang-java";
import { cppLanguage, cpp } from "@codemirror/lang-cpp";
import { goLanguage, go } from "@codemirror/lang-go";
import "./styles.scss";
import { message, Select } from "antd";
import { LanguageSupport } from "@codemirror/language";
import { ProgrammingLanguageOptions } from "@/utils/SelectOptions";
import { ExecutionResults, SubmissionResults } from "@/app/services/execute";
import {language} from "@codemirror/language"
// import {htmlLanguage, html} from "@codemirror/lang-html"


interface CollaborativeEditorProps {
  user: string;
  collaborationId: string;
  language: string;
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
  }
  executingState: {
    executing: boolean;
    id: number;
  }
  submittingState: {
    submitting: boolean;
    id: number;
  }
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
    let [selectedLanguage, setSelectedLanguage] = useState<string>("Python");
    let sessionEndNotified = false;

    const languageConf = new Compartment();

    const codeChangeListener = EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        props.onCodeChange(update.state.doc.toString());
      }
    });

    const autoLanguage = EditorState.transactionExtender.of(tr => {
      if (!tr.docChanged) return null

      let storedLanguageLabel = localStorage.getItem("languageLabel") as string;
      let newLanguage: Extension = python();
      let languageType;
      let languageLabel: string = "";

      if (storedLanguageLabel.toLowerCase() == "python") {
        newLanguage = python();
        languageLabel = "Python";
        languageType = pythonLanguage;
      } else if (storedLanguageLabel.toLowerCase() == "javascript") {
        newLanguage = javascript();
        languageLabel = "JavaScript";
        languageType = javascriptLanguage;
      }

      const stateLanguage = tr.startState.facet(language);
      if (languageType == stateLanguage) return null;

      setSelectedLanguage(languageLabel);

      return {
        effects: languageConf.reconfigure(newLanguage),
      };
    })
    
    
    // Referenced: https://codemirror.net/examples/config/#dynamic-configuration
    // const autoLanguage = EditorState.transactionExtender.of((tr) => {
    //   if (!tr.docChanged) return null;

    //   const snippet = tr.newDoc.sliceString(0, 100);

    //   // Handle code change
    //   props.onCodeChange(tr.newDoc.toString());

    //   // Test for various language
    //   const docIsPython = /^\s*(def|class)\s/.test(snippet);
    //   const docIsJava = /^\s*(class|public\s+static\s+void\s+main)\s/.test(
    //     snippet
    //   ); // Java has some problems
    //   const docIsCpp = /^\s*(#include|namespace|int\s+main)\s/.test(snippet); // Yet to test c++
    //   const docIsGo = /^(package|import|func|type|var|const)\s/.test(snippet);

    //   let newLanguage;
    //   let languageType;
    //   let languageLabel;

    //   if (docIsPython) {
    //     newLanguage = python();
    //     languageLabel = "Python";
    //     languageType = pythonLanguage;
    //   } else if (docIsJava) {
    //     newLanguage = java();
    //     languageLabel = "Java";
    //     languageType = javaLanguage;
    //   } else if (docIsGo) {
    //     newLanguage = go();
    //     languageLabel = "Go";
    //     languageType = goLanguage;
    //   } else if (docIsCpp) {
    //     newLanguage = cpp();
    //     languageLabel = "C++";
    //     languageType = cppLanguage;
    //   } else {
    //     newLanguage = javascript(); // Default to JavaScript
    //     languageLabel = "JavaScript";
    //     languageType = javascriptLanguage;
    //   }

    //   const stateLanguage = tr.startState.facet(language);
    //   if (languageType == stateLanguage) return null;

    //   setSelectedLanguage(languageLabel);

    //   return {
    //     effects: languageConf.reconfigure(newLanguage),
    //   };
    // });

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

    let latestExecutionId: number = (new Date(0)).getTime();
    let latestSubmissionId: number = (new Date(0)).getTime();
    let latestExecutingId: number = (new Date(0)).getTime();
    let latestSubmittingId: number = (new Date(0)).getTime();

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

    let sessionEndTimeout: any;


    let view = new EditorView();

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

            if (
              state && 
              state.submissionResultsState &&
              state.submissionResultsState.id !== latestSubmissionId
            ) {
              latestSubmissionId = state.submissionResultsState.id;
              props.updateSubmissionResults(state.submissionResultsState.submissionResults);
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
              props.updateExecutionResults(state.executionResultsState.executionResults);
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
                  content: `${
                    props.matchedUser ?? "Peer"
                  } is saving code...`,
                });
              }
            }
          });
      });

      let state = EditorState.create();

      // if (selectedLanguage == "python" || selectedLanguage == "Python") {
        state = EditorState.create({
          doc: ytext.toString(),
          extensions: [
            basicSetup,
            languageConf.of(python()),
            // languageConf.of(javascript()),
            autoLanguage,
            // changeLanguage,
            yCollab(ytext, provider.awareness, { undoManager }),
            keymap.of([indentWithTab]),
            codeChangeListener,
          ],
        });
      // }
      // else if (selectedLanguage == "java" || selectedLanguage == "Java") {
      //   state = EditorState.create({
      //     doc: ytext.toString(),
      //     extensions: [
      //       basicSetup,
      //       languageConf.of(java()),
      //       // languageConf.of(javascript()),
      //       // autoLanguage,
      //       // changeLanguage,
      //       yCollab(ytext, provider.awareness, { undoManager }),
      //       keymap.of([indentWithTab]),
      //       codeChangeListener,
      //     ],
      //   });
      // }

      view = new EditorView({
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
            defaultValue={selectedLanguage}
            options={ProgrammingLanguageOptions}
            onSelect={(val) => {
              setSelectedLanguage(val);
              localStorage.setItem("languageLabel", val);
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

const EXECUTION_SERVICE_URL = process.env.NEXT_PUBLIC_EXECUTION_SERVICE_URL;

export interface TestData {
  input: string;
  expected: string;
}

export interface TestResult {
  input: string;
  expected: string;
  actual: string;
  passed: boolean;
  error: string;
}

export interface AllTests {
  visibleTests: TestData[];
  hiddenTests: TestData[];
}

export type Test = TestResult | TestData;

export const isTestResult = (test: Test): test is TestResult => {
  return "actual" in test && "passed" in test && "error" in test;
};

function parseTestcases(tests: Test[]) {
  return tests.reduce((accum: string, t: Test) => {
    return accum + t.input + "\n" + t.expected + "\n";
  }, `${tests.length}\n`);
}

export interface GeneralTestResults {
  passed: number;
  total: number;
}

export interface SubmissionHiddenTestResultsAndStatus {
  hiddenTestResults: GeneralTestResults;
  status: string;
}

export interface SubmissionResults
  extends SubmissionHiddenTestResultsAndStatus {
  visibleTestResults: TestResult[];
}

export interface ExecutionResults {
  visibleTestResults: TestResult[];
  customTestResults: TestResult[];
}

export interface Code {
  code: string;
  language: string;
  customTestCases: string;
}

export interface Submission {
  title: string;
  code: string;
  language: string;
  user: string;
  matchedUser: string;
  matchedTopics: string[];
  questionDifficulty: string;
  questionTopics: string[];
}

export const GetVisibleTests = async (
  questionDocRefId: string
): Promise<TestData[]> => {
    const response = await fetch(
        `${EXECUTION_SERVICE_URL}tests/${questionDocRefId}`,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        }
    );

    if (response.ok) {
        return response.json();
    } else {
        throw new Error(
            `Error fetching test cases: ${await response.text()}`
        );
    }
};

export const ExecuteVisibleAndCustomTests = async (
  questionDocRefId: string,
  code: Code
): Promise<ExecutionResults> => {
    const response = await fetch(
        `${EXECUTION_SERVICE_URL}tests/${questionDocRefId}/execute`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(code),
        }
    );

    if (response.ok) {
        return response.json();
    } else {
        throw new Error(
            `Error executing code: ${response.status} ${response.statusText}`
        );
    }
};

export const ExecuteVisibleAndHiddenTestsAndSubmit = async (
  questionDocRefId: string,
  collaboration: Submission
): Promise<SubmissionResults> => {
    const response = await fetch(
        `${EXECUTION_SERVICE_URL}tests/${questionDocRefId}/submit`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(collaboration),
        }
    );

    if (response.ok) {
        return response.json();
    } else {
        throw new Error(
            `Error submitting code: ${response.status} ${response.statusText}`
        );
    }
};

export const CreateTestcases = async (
  questionDocRefId: string,
  questionTitle: string,
  visibleTestcases: TestData[],
  hiddenTestcases: TestData[]
): Promise<any> => {
  const response = await fetch(`${EXECUTION_SERVICE_URL}tests`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      questionDocRefId: questionDocRefId,
      questionTitle: questionTitle,
      visibleTestCases: parseTestcases(visibleTestcases),
      hiddenTestCases: parseTestcases(hiddenTestcases),
    }),
  });

  if (response.ok) {
    return response.json();
  } else {
    throw new Error(
      `Error creating testcases: ${response.status} ${response.statusText}`
    );
  }
};

export const ReadAllTestcases = async (
  questionDocRefId: string
): Promise<AllTests> => {
  const response = await fetch(
    `${EXECUTION_SERVICE_URL}tests/${questionDocRefId}/readall`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  if (response.ok) {
    return response.json();
  } else {
    throw new Error(
      `Error reading testcases: ${response.status} ${response.statusText}`
    );
  }
};

export const UpdateTestcases = async (
  questionDocRefId: string,
  visibleTestcases: TestData[],
  hiddenTestcases: TestData[]
): Promise<any> => {
  const response = await fetch(
    `${EXECUTION_SERVICE_URL}tests/${questionDocRefId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        visibleTestCases: parseTestcases(visibleTestcases),
        hiddenTestCases: parseTestcases(hiddenTestcases),
      }),
    }
  );

  if (response.ok) {
    return response.json();
  } else {
    throw new Error(
      `Error updating testcases: ${response.status} ${response.statusText}`
    );
  }
};

export const DeleteTestcases = async (
  questionDocRefId: string
): Promise<void> => {
  const response = await fetch(
    `${EXECUTION_SERVICE_URL}tests/${questionDocRefId}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  if (!response.ok) {
    throw new Error(
      `Error deleting testcases: ${response.status} ${response.statusText}`
    );
  }
};

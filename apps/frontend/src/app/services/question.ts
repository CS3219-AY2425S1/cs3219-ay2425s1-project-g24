import { getToken } from "./login-store";

const QUESTION_SERVICE_URL = process.env.NEXT_PUBLIC_QUESTION_SERVICE_URL;

export interface Question {
  id: number;
  docRefId: string;
  title: string;
  description?: string;
  categories: string[]; // enum[]
  complexity: string; // enum
  assets?: string[];
  createdAt?: Date; // or string
  updatedAt?: Date; // or string
  deletedAt?: Date;
  testCases?: string[];
}

export interface QuestionListResponse {
  totalCount: number;
  totalPages: number;
  currentPage: number;
  limit: number;
  hasNextPage: boolean;
  questions: Question[];
}

export interface NewQuestion {
  title: string;
  description: string;
  categories: string[];
  complexity: string;
}

// GET request to fetch all the questions (TODO: Ben --> Fetch with filtering/sorting etc)
export const GetQuestions = async (
  currentPage?: number,
  limit?: number,
  sortBy?: string,
  difficulties?: string[],
  categories?: string[],
  title?: string
): Promise<QuestionListResponse> => {
  let query_url = `${QUESTION_SERVICE_URL}questions`;
  let query_params = "";

  if (currentPage) {
    query_params += `?offset=${(currentPage - 1) * 10}`;
  }

  if (limit) {
    query_params += `${query_params.length > 0 ? "&" : "?"}limit=${limit}`;
  }

  if (sortBy) {
    const [field, order] = sortBy.split(" ");
    query_params += `${
      query_params.length > 0 ? "&" : "?"
    }sortField=${field}&sortValue=${order}`;
  }

  if (difficulties && difficulties.length > 0) {
    const value = difficulties.join(","); // Change them from ["easy", "medium"] to "easy,medium"
    query_params += `${query_params.length > 0 ? "&" : "?"}complexity=${value}`;
  }

  if (categories && categories.length > 0) {
    const value = categories.join(",");
    query_params += `${query_params.length > 0 ? "&" : "?"}categories=${value}`;
  }

  if (title && title != "") {
    const urlEncodedTitle = encodeURIComponent(title);
    query_params += `${
      query_params.length > 0 ? "&" : "?"
    }title=${urlEncodedTitle}`;
  }

  query_url += query_params;
  const response = await fetch(query_url, {
    method: "GET",
    headers: {
      'Authorization': `Bearer ${getToken()}`,
    }
  });
  const data = await response.json();
  return data;
};

// Get single question
export const GetSingleQuestion = async (docRef: string): Promise<Question> => {
  const response = await fetch(
    `${QUESTION_SERVICE_URL}questions/${docRef}`,
    {
      method: "GET",
      headers: {
        'Authorization': `Bearer ${getToken()}`,
      }
    }
  );
  const data = await response.json();
  return data;
};

// Upload single question (TODO: Sean)
export const CreateQuestion = async (
  question: NewQuestion
): Promise<Question> => {
  const response = await fetch(`${QUESTION_SERVICE_URL}questions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      'Authorization': `Bearer ${getToken()}`,
    },
    body: JSON.stringify(question),
  });

  if (response.status === 200) {
    return response.json();
  } else {
    throw new Error(
      `Error creating question: ${response.status} ${response.statusText}`
    );
  }
};

export const EditQuestion = async (
  question: NewQuestion,
  docRefId: string
): Promise<Question> => {
  const response = await fetch(
    `${QUESTION_SERVICE_URL}questions/${docRefId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${getToken()}`,
      },
      body: JSON.stringify(question),
    }
  );

  if (response.status === 200) {
    return response.json();
  } else {
    throw new Error(
      `Error updating question: ${response.status} ${response.statusText}`
    );
  }
};

// Delete single question (TODO: Ryan)
export async function DeleteQuestion(docRef: String): Promise<void> {
  const res = await fetch(
    `${QUESTION_SERVICE_URL}questions/${docRef}`,
    {
      method: "DELETE",
      headers: {
        'Authorization': `Bearer ${getToken()}`,
      },
    }
  );
  // error handling later
  // if (!res.ok) {
  //   throw new Error(`Delete request responded with ${res.status}: ${res.body}`)
  // }
}

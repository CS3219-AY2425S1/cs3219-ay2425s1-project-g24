import { CreateQuestion, DeleteQuestion, GetQuestions, GetSingleQuestion, Question } from "@/app/services/question"

const NEXT_PUBLIC_QUESTION_SERVICE_URL = process.env.NEXT_PUBLIC_QUESTION_SERVICE_URL

const QUESTIONS = [
    {
      "id": 12345,
      "docRefId": "asdfDocRef",
      "title": "Asdf",
      "description": "Asdf description",
      "categories": ["Arrays", "Algorithms"],
      "complexity": "hard"
    },
    {
      "id": 12346,
      "docRefId": "qwerDocRef",
      "title": "Qwer",
      "description": "Qwer description",
      "categories": ["Strings", "Data Structures"],
      "complexity": "medium"
    },
    {
      "id": 12347,
      "docRefId": "zxcvDocRef",
      "title": "Zxcv",
      "description": "Zxcv description",
      "categories": ["Graphs", "Algorithms"],
      "complexity": "easy"
    },
    {
      "id": 12348,
      "docRefId": "tyuiDocRef",
      "title": "Tyui",
      "description": "Tyui description",
      "categories": ["Trees", "Recursion"],
      "complexity": "hard"
    },
    {
      "id": 12349,
      "docRefId": "ghjkDocRef",
      "title": "Ghjk",
      "description": "Ghjk description",
      "categories": ["Dynamic Programming", "Math"],
      "complexity": "medium"
    },
    {
      "id": 12350,
      "docRefId": "bnmlDocRef",
      "title": "Bnml",
      "description": "Bnml description",
      "categories": ["Sorting", "Searching"],
      "complexity": "easy"
    },
    {
      "id": 12351,
      "docRefId": "poiuDocRef",
      "title": "Poiu",
      "description": "Poiu description",
      "categories": ["Bit Manipulation", "Algorithms"],
      "complexity": "hard"
    },
    {
      "id": 12352,
      "docRefId": "lkjhDocRef",
      "title": "Lkjh",
      "description": "Lkjh description",
      "categories": ["Greedy", "Data Structures"],
      "complexity": "medium"
    },
    {
      "id": 12353,
      "docRefId": "mnbvDocRef",
      "title": "Mnbv",
      "description": "Mnbv description",
      "categories": ["Backtracking", "Recursion"],
      "complexity": "easy"
    },
    {
      "id": 12354,
      "docRefId": "vcxzDocRef",
      "title": "Vcxz",
      "description": "Vcxz description",
      "categories": ["Graphs", "Dynamic Programming"],
      "complexity": "hard"
    }
]

jest.mock("@/app/services/login-store", () => {
    return {
        __esModule: true,
        getToken: jest.fn(() => TOKEN)
    };
})

function createMockResponse(obj: any): Response {
    // @ts-ignore don't need the whole response
    return { 
        json: () => Promise.resolve(obj),
        ok: true,
        status: 200
    }
}

const TOKEN = "mockjwttoken"

describe("GetQuestions", () => {
    beforeEach(() => {
        global.fetch = jest.fn().mockResolvedValue({
            async json() {
                return QUESTIONS
            }
        });
    })
    
    it("gets all questions on the first page with () call", async () => {
        
        const res = await GetQuestions()

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions`, {
            headers: {
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "GET",
        }]])
        expect(res).toStrictEqual(QUESTIONS)

    });
    
    it("formats (page=2) params correctly", async () => {
        
        const res = await GetQuestions(2)

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions?offset=10`, {
            headers: {
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "GET",
        }]])
    });
    
    it("formats (limit=3) params correctly", async () => {
        
        await GetQuestions(undefined, 3)

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions?limit=3`, {
            headers: {
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "GET",
        }]])
    });
    
    it("formats (difficulty asc) params correctly", async () => {
        
        await GetQuestions(undefined, undefined, "difficulty asc")

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions?sortField=difficulty&sortValue=asc`, {
            headers: {
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "GET",
        }]])
    });
    
    it("formats ([\"easy\", \"hard\"]) params correctly", async () => {
        
        await GetQuestions(undefined, undefined, undefined, ["easy", "hard"])

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions?complexity=easy,hard`, {
            headers: {
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "GET",
        }]])
    });
    
    it("formats cat params correctly", async () => {
        
        await GetQuestions(undefined, undefined, undefined, undefined, ["CatA", "CatB"])

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[
            `${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions?categories=CatA,CatB`, 
            {
                headers: {
                    "Authorization": `Bearer ${TOKEN}`,
                },
                method: "GET",
            }
        ]])
    });
    
    it("formats title params correctly", async () => {
        
        await GetQuestions(undefined, undefined, undefined, undefined, undefined, "The Title Name")

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[
            `${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions?title=The%20Title%20Name`, 
            {
                headers: {
                    "Authorization": `Bearer ${TOKEN}`,
                },
                method: "GET",
            }
        ]])
    });
})


describe("GetSingleQuestion", () => {
    const DOCREF = "mockdocref";
    beforeEach(() => {
        global.fetch = jest.fn().mockResolvedValue({
            ok: true, // Ensure `ok` is true to hit the success branch
            async json() {
                return QUESTIONS[0]
            },
            text: () => Promise.resolve('mocked response'),
        });
    });

    it("gets a question by docref", async () => {
        const res = await GetSingleQuestion(DOCREF);
        
        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions/${DOCREF}`, {
            headers: {
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "GET",
        }]])
        expect(res).toStrictEqual(QUESTIONS[0]);
    })
})

describe("CreateQuestion", () => {
    it("uploads a question", async () => {
        // grabs a subset of QUESTIONS[0]
        const newQuestion = (({title, description, categories, complexity}) => ({title, description, categories, complexity}))(QUESTIONS[0])
        const createdQuestion = QUESTIONS[0];

        global.fetch = jest.fn().mockResolvedValue({
            status: 200,
            statusText: "OK",
            ok: true, // Ensure `ok` is true to hit the success branch
            async json() {
                return createdQuestion
            }
        });

        const res = await CreateQuestion(newQuestion);
        
        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions`, {
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "POST",
            body: JSON.stringify(newQuestion)
        }]])
        expect(res).toStrictEqual(createdQuestion);
    })

    it("fails uploading question", async () => {
        // grabs a subset of QUESTIONS[0]
        const newQuestion = (({title, description, categories, complexity}) => ({title, description, categories, complexity}))(QUESTIONS[0])

        global.fetch = jest.fn().mockResolvedValue({
            status: 400,
            statusText: "Not Found",
            data: "Question title already exists"
        })

        const res = CreateQuestion(newQuestion);
        
        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[`${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions`, {
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${TOKEN}`,
            },
            method: "POST",
            body: JSON.stringify(newQuestion)
        }]])
        await expect(res).rejects.toThrow("Error creating question: 400 Not Found")
    })


})


describe("DeleteQuestion", () => {
    const DOCREF = "mockdocref";
    beforeEach(() => {
        global.fetch = jest.fn().mockResolvedValue({
            status: 200,
            statusText: "OK",
            data: `Question with ID ${DOCREF} deleted successfully`
        });
    });

    it("deletes successfully", async () => {
        const shouldbeNothing = await DeleteQuestion(DOCREF);

        expect(jest.mocked(fetch).mock.calls).toStrictEqual([[
            `${NEXT_PUBLIC_QUESTION_SERVICE_URL}questions/${DOCREF}`, 
            {
                headers: {
                    "Authorization": `Bearer ${TOKEN}`,
                },
                method: "DELETE",
            }
        ]])
        expect(shouldbeNothing).toBeUndefined();
    })
})

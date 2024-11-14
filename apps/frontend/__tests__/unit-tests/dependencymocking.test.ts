import { GetQuestions } from "@/app/services/question"
import { getToken } from "@/app/services/login-store";

const TOKEN = 'mocktoken';

jest.mock("@/app/services/login-store", () => {
    return {
        __esModule: true,
        getToken: jest.fn(() => TOKEN)
    };
})

beforeEach(() => {
    global.fetch = jest.fn().mockResolvedValue({
        async json() {
            return {}
        }
    });
})

describe("mock", () => {
    
    it("mocks correctly", async () => {
        await GetQuestions()
        expect(jest.mocked(getToken).mock.calls).toEqual([[]])
        expect(jest.mocked(fetch).mock.calls[0][1]).toEqual({
            "headers": {
                "Authorization": `Bearer ${TOKEN}`,
            },
            "method": "GET",
        })
    });
    
})
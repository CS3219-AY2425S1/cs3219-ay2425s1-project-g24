import { Actions, Browser, Builder, By, Capabilities, Key, until, WebDriver } from "selenium-webdriver"

import {Options as ChromeOptions} from "selenium-webdriver/chrome"
import {Options as EdgeOptions} from "selenium-webdriver/edge"
import {Options as FirefoxOptions} from "selenium-webdriver/firefox"

const URL = 'http://localhost:3000/';
const ETERNAL_JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjk5OTk5OTk5OTk5fQ.Z4_FVGQ5lIcouP3m4YLMr6pGMF17IJFfo2yOTiN58DY"

const CHROME_OPTIONS = new ChromeOptions()
    .addArguments("--headless=new") as ChromeOptions; // uncomment locally to see the steps in action
const EDGE_OPTIONS = new EdgeOptions()
    .setBinaryPath("/opt/hostedtoolcache/msedge/stable/x64/msedge") // need to point to the correct path
    .addArguments("--headless=new") as EdgeOptions;

const FIREFOX_OPTIONS = new FirefoxOptions()
    .addArguments("--headless") as FirefoxOptions;

const builder = new Builder()
    .setChromeOptions(CHROME_OPTIONS)
    .setEdgeOptions(EDGE_OPTIONS)
    .setFirefoxOptions(FIREFOX_OPTIONS)

describe.each([Browser.CHROME, Browser.EDGE, Browser.FIREFOX])("%s driver test", (browser) => {
    let driver: WebDriver;
    beforeAll(() => {
        const cap = new Capabilities().setBrowserName(browser)
        builder.withCapabilities(cap);
    })

    beforeEach(async () => {
        console.log(browser + ": building...");
        driver = await builder.build();
        console.log(browser + ": built");
    }, 20000)

    afterEach(async () => {
        if (driver) {
            await driver.quit();
        }
    })

    describe("webdriver installed correctly", () => {
        it("does google search", async () => {
            await driver.get('http://www.google.com');
            await driver.findElement(By.name('q')).sendKeys('webdriver', Key.RETURN);
            await driver.wait(until.titleIs('webdriver - Google Search'), 1000);
        }, 10000);

        it.skip("does another google search", async () => {
            await driver.get('http://www.google.com');
            await driver.findElement(By.name('q')).sendKeys('webdriver', Key.RETURN);
            await driver.wait(until.titleIs('webdriver - Google Search'), 1000);
        }, 10000);
    });
    
    describe("browser-test", () => {
        it("accesses and login to peerprep", async () => {
            await driver.get(URL);
            await driver.wait(until.urlIs(`${URL}login`));
            
            const [email, password] = await driver.findElements(By.css("input"))
            const submit = await driver.findElement(By.css("button[type=\"submit\"]"))

            await email.sendKeys("admin@gmail.com");
            await password.sendKeys("admin");

            await submit.click();
            await driver.wait(until.urlIs(`${URL}`));

            const slogan1 = await driver.findElement(By.xpath("/html/body/div[1]/main/div/div[1]/div[2]/span[1]")).then(ele => ele.getText())
            const slogan2 = await driver.findElement(By.xpath("/html/body/div[1]/main/div/div[1]/div[2]/span[2]")).then(ele => ele.getText())

            expect(slogan1).toBe("A better way to prepare for coding interviews with");
            expect(slogan2).toBe("peers");
        }, 10000);
    })
}, 60000)





// shamelessly stolen from https://github.com/daleroy1/PoEWikiBot
// and hacked apart for use here

const puppeteer = require('puppeteer');
const wikiURL = "https://pathofexile.gamepedia.com/";
const wikiDiv = ".infobox-page-container > .item-box";
const wikiInfoDiv = ".infocard"
const wikiInvalidPage = ".noarticletext"

var browser;

(async () => {
    browser = await puppeteer.launch({
        ignoreHTTPSError: true,
        headless: true,
        handleSIGHUP: true,
        args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage']
    });

    res = await getItemImage(process.argv[2])
    if (res.success = true && res.screenshot != false) {
        console.log(res.screenshot.toString('base64'))
        process.exit()
    } else {
        process.exit(1)
    }
})();



async function getItemImage(name) {
    let itemUrlPart = convertToUrlString(titleCase(name));
    var url = wikiURL + itemUrlPart;

    //console.time('getPage')
    const page = await browser.newPage();
    //Disabling Javascript adds 100% increased performance
    await page.setJavaScriptEnabled(false)
    var output = {
        screenshot: false,
        success: false
    }

    //Set a tall page so the image isn't covered by popups 
    await page.setViewport({ 'width': 2500, 'height': 2500 });

    try {
        //played around with a few different waitUntils.  This one seemed the quickest.
        //If you don't disable Javascript on the PoE Wiki site, removing this parameter makes it hang
        await page.goto(url, { waitUntil: 'load' });
    } catch (error) {
        console.error(`"${error.message}" "${url}"`);
    }

    var invalidPage = await page.$(wikiInvalidPage);
    //if we have a invalid page, lets exit
    if (invalidPage != null) {
        return output;
    }

    var infoBox = await page.$(wikiInfoDiv);
    if (infoBox != null) {
        try {
            output.screenshot = await infoBox.screenshot();
            output.success = true;
        } catch (error) {
            output.success = true;
        }
        return output;
    }

    //if we have a div for the item, screenshot it.
    //If not, just return the page without the screenshot
    const div = await page.$(wikiDiv);
    if (div != null) {
        try {
            output.screenshot = await div.screenshot();
            output.success = true;
        } catch (error) {
            output.success = true;
        }
    } else {
        output.success = true;
    }

    await page.close();
    return output;
}

function convertToUrlString(name) {
    return name.replace(new RegExp(" ", "g"), "_");
}

function titleCase(str) {
    let excludedWords = ["of", "and", "the", "to", "at", "for"];
    let words = str.split(" ");
    for (var i in words) {
        if ((i == 0) || !(excludedWords.includes(words[i].toLowerCase()))) {
            words[i] = words[i][0].toUpperCase() + words[i].slice(1, words[i].length);
        } else {
            continue;
        }
    }
    return words.join(" ");
};
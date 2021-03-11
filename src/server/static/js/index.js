let LEDS = []
let Modes = []

async function getLeds() {
    let resp = await fetch("/leds")
    if (resp.ok) {
        LEDS = await resp.json()
        console.info("INFO", LEDS)
    } else {
        console.error("ERROR", resp.status)
    }
}

async function getModes() {
    let resp = await fetch("/modes")
    if (resp.ok) {
        Modes = await resp.json()
        console.info("INFO", Modes)
    } else {
        console.error("ERROR", resp.status)
    }
}

async function sendLeds() {
    let resp = await fetch("/leds", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(LEDS)
    })
}

function changeBrightness(place) {
    let selector = `#${place.id} input[type='range']`
    let br = document.querySelector(selector).value
    for (i = 0; i < LEDS.length; i++) {
        if (LEDS[i].Place == place.id) {
            LEDS[i].Brightness = br / 100
        }
    }
    sendLeds()
}

function changeColor(place) {
    let selector = `#${place.id} input[type='color']`
    let color = document.querySelector(selector).value
    for (i = 0; i < LEDS.length; i++) {
        if (LEDS[i].Place == place.id) {
            LEDS[i].Color = color
        }
    }
    sendLeds()
}

function changeMode(place) {
    let selector = `#${place.id} select`
    let mode = document.querySelector(selector).value
    for (i = 0; i < LEDS.length; i++) {
        if (LEDS[i].Place == place.id) {
            LEDS[i].Mode = mode
        }
    }
    sendLeds()
}

function renderLEDSInterface() {

    let out = ""
    for (i = 0; i < LEDS.length; i++) {
        l = LEDS[i]

        s = `
        <div class="led" id="${l.Place}">
            <h1>${l.Place}</h1>
            <div class="form-group">
        `

        // render modes
        s += `<select class="form-control form-control-lg" oninput="changeMode(${l.Place})">`
        for (j = 0; j < Modes.length; j++) {
            if (LEDS[i].Mode == Modes[j]) {
                // current active mode
                s += `
                <option selected value="${Modes[j]}">${Modes[j]}</option>
                `
            } else {
                s += `
                <option value="${Modes[j]}">${Modes[j]}</option>
                `
            }
        }
        s += `</select>`

        // brightness and color
        s += `
            <p>
            <label>Brightness
            <input type="range" class="form-control-range"
                oninput="changeBrightness(${l.Place});" value="${l.Brightness * 100}">
            </label>
            </p>
            <p>
            <label>Color
            <input type="color" value="${l.Color}"
                oninput="changeColor(${l.Place})">
            </label>
            </p>
        `
        // end rendering string
        s += `</div></div>`

        out += s
    }
    document.querySelector("#leds").innerHTML = out
}

function updateData() {
    for (i = 0; i < LEDS.length; i++) {
        // let el = document.querySelector(`#${LEDS[i].Place}`)
        let br = document.querySelector(`#${LEDS[i].Place} input[type='range']`)
        let co = document.querySelector(`#${LEDS[i].Place} input[type='color']`)
        let se = document.querySelector(`#${LEDS[i].Place} select`)
        br.value = LEDS[i].Brightness * 100
        co.value = LEDS[i].Color
        se.value = LEDS[i].Mode
    }
}

getModes().then( () => {
    getLeds().then(renderLEDSInterface)
})

setInterval(() => {
    getLeds().then(updateData)
}, 100)
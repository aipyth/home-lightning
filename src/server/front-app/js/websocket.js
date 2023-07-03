const max_connect_retries = 5
const connection = {
    retries_left: max_connect_retries,
    connected: false,
    socket: undefined,

    data: {
        modes: [],
        places: [],
    },

    requestModes() { this.socket.send("get-modes") },
    requestPlaces() { this.socket.send("get-places") },
    createMode(name) { this.socket.send(`create-mode;${name}`) },
    createPlace(name) { this.socket.send(`create-place;${name}`) },
    removeMode(name) { this.socket.send(`remove-mode;${name}`) },
    removePlace(name) { this.socket.send(`remove-place;${name}`) },
    updatePlace(name, mode, color, brightness) {
        this.socket.send(`update-place;${name};${mode};${color};${brightness}`)
    },


    connect() {
        const websocketProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const websocketHost = window.location.host
        const websocketPath = 'socket/web-interface'
        const websocketUrl = `${websocketProtocol}//${websocketHost}/${websocketPath}`

        this.socket = new WebSocket(websocketUrl)

        this.socket.onopen = () => {
            console.debug('Socket connection opened')
            this.connected = true
            this.retries_left = max_connect_retries
        }

        this.socket.onclose = () => {
            this.connected = false
            console.debug('Socket connection closed')
            if (this.retries_left !== 0) {
                this.connect()
                this.retries_left--
            } else {
                console.log("No retries left for connection to server")
            }
        }

        this.socket.onmessage = (ev) => {
            const data = ev.data
            console.log("from server:", data)

            if (data[0] === '!') {
                // TODO make notifications
                console.log(data[0])
                return
            }
            const args = data.split(';')
            // first argument identifies the query
            switch (args[0]) {
                case "modes":
                    this.data.modes = args.slice(1, -1)
                    break
                case "places":
                    this.data.places = args.slice(1, -1)
                    break
                default:
                    // some place got updated info
                    const place = args[1]
                    if (!this.data.places.includes(place)) {
                        console.log('Wrong query from server.', args)
                        return
                    }
                    const mode = args[2]
                    const color = args[3]
                    const brightness = parseFloat(args[4])
                    this.data[place] = { mode, color, brightness }
            }
        }
    },
}

connection.connect()

export default connection
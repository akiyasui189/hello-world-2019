const app = require('express')()
const http = require('http').createServer(app);
const io = require('socket.io').listen(http);
const request = require('request')

app.get('/', (req, res) => res.sendFile(__dirname + '/index.html'))

io.on('connection', function(socket) {
    // start connection
    console.log('user connected');
    // receive message
    socket.on('chat message', function(msg) {
        console.log('receive message: ' + msg);
        // call API
        request.post({
            url: "http://localhost:8080/helloworld",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                message: msg
            }),
            json: true
        }, function (error, response, body) {
            console.log('response: ' + JSON.stringify(response))
            console.log('response body: ' + JSON.stringify(body))
            messageId = body.messageId;
            registeredAt = body.registeredAt;
            var data = {
                id: messageId,
                message: msg,
                registeredAt: registeredAt
            };
            // send message
            io.emit('chat message', data);
        });
    });
    // disconnect
    socket.on('disconnect', function() {
        console.log('user disconnected');
    });
});

http.listen(3000, () => console.log('FE app listening on port 3000'))

const app = require('express')()
const http = require('http').createServer(app);
const io = require('socket.io').listen(http);

app.get('/', (req, res) => res.sendFile(__dirname + '/index.html'))

io.on('connection', function(socket) {
    // start connection
    console.log('a user connected');
    // receive message
    socket.on('chat message', function(msg) {
        console.log('message: ' + msg);
        io.emit('chat message', msg);
    });
    // disconnect
    socket.on('disconnect', function() {
        console.log('user disconnected');
    });
});

http.listen(3000, () => console.log('Example app listening on port 3000!'))

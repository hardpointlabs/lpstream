import { Encoder, Decoder } from '@hardpointlabs/length-prefixed-stream';
import * as net from 'node:net';

const server = net.createServer((socket) => {
  const encoder = new Encoder();
  const decoder = new Decoder();
  console.log('Client connected');
  socket
  .pipe(encoder)
  .pipe(decoder)
  .pipe(socket);
});

server.listen(8124, () => {
  console.log('TCP echo server listening on port 8124');
});

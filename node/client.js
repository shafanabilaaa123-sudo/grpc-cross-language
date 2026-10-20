const path = require('path');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

// Path ke file proto utama (berada satu level di atas folder node)
const PROTO_PATH = path.join(__dirname, '../user.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const pb = protoDescriptor.pb;

function main() {
    const client = new pb.UserService(
        'localhost:50051',
        grpc.credentials.createInsecure()
    );

    const userId = 1;
    console.log(`Mengirim request untuk mengambil User ID: ${userId}...`);

    client.GetUser({ id: userId }, (err, response) => {
        if (err) {
            console.error('Error saat memanggil gRPC:', err);
            return;
        }
        console.log('Respon dari Server Go:');
        console.log(response);
    });
}

main();
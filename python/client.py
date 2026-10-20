import grpc
import user_pb2
import user_pb2_grpc

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = user_pb2_grpc.UserServiceStub(channel)
    
    print("Mengirim request Python ke Go Server...")
    response = stub.GetUser(user_pb2.UserRequest(id=1))
    
    print("Respon dari Server Go:")
    print(f"ID: {response.id}, Name: {response.name}, Phone: {response.phone}")

if __name__ == '__main__':
    run()
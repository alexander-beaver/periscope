import os
from google.protobuf.json_format import MessageToJson
from concurrent import futures
from io import BytesIO
import json
import grpc
from pypdf import PdfReader
from plugin_pb2 import LoadResponse, PluginCapability, TransformRequest, TransformEntry,TransformEntryType, TerminateResponse, TransformResponse
import requests
import sys
import threading
from socketserver import ThreadingMixIn
import random

from plugin_pb2_grpc import PluginServicer, TransformerServicer, add_PluginServicer_to_server, add_TransformerServicer_to_server
from http.server import BaseHTTPRequestHandler, HTTPServer

import time



def terminate():
    print("Terminating")
    time.sleep(5)
    os._exit(1)

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # Extracting the 'id' from the URL path
        id_str = self.path.split('/')[-1]
        print("PDF Path",self.path)
        try:
            # Convert the id to an integer
            print("GETTING ID STR",id_str)
            id = int(id_str)
            print("ID",id)
            # Access the entry from PDFTransformer class
            entry = self.server.transformer.entries[id]
            print("Entry",entry)
            if entry != None:
                # Respond with the entry contents
                print("Responding")
                print("Entry encoded",entry)
                self.send_response(200)
                self.send_header("Content-type", "application/json")
                self.end_headers()
                self.wfile.write(entry.encode(encoding='utf_8'))
                print("Wrote")
                # Delete cache file
                os.remove("{}.cache".format(id))
                
                
            else:
                '''
                if os.path.exists("{}.cache".format(id)):
                    with open("{}.cache".format(id), "r") as f:
                        entry = json.load(f)
                        print("Entry",entry)
                        self.send_response(200)
                        self.send_header("Content-type", "application/json")
                        self.end_headers()
                        self.wfile.write(json.dumps(entry).encode(encoding='utf_8'))
                    os.remove("{}.cache".format(id))
                else:
                    print("Entry not found")
                    # Entry not found'''
                
                self.send_response(404)
                self.end_headers()
                entry = TransformEntry(type=TransformEntryType.NOINDEX, contents="Entry not found")
                self.wfile.write(MessageToJson(entry).encode(encoding='utf_8'))
                    
                    
        except ValueError:
            print("Invalid ID")
            self.send_response(400)
            self.end_headers()
            entry = TransformEntry(type=TransformEntryType.NOINDEX, contents="Invalid ID")
            self.wfile.write(MessageToJson(entry).encode(encoding='utf_8'))
        except Exception as e:
            print("Exception",e)
            self.send_response(500)
            self.end_headers()
            entry = TransformEntry(type=TransformEntryType.NOINDEX, contents=e)
            self.wfile.write(MessageToJson(entry).encode(encoding='utf_8'))

class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in a separate thread."""

class PDFTransformer(TransformerServicer, PluginServicer):
    entries = {}
    grpcPort = 0000
    httpPort = 0000
    def start_http_server(self):
        server_address = ('', self.httpPort)
        httpd = ThreadedHTTPServer(server_address, RequestHandler)
        httpd.transformer = self  # Passing the transformer reference to the handler
        httpd.serve_forever()

    def Load(self, request, context):

        res  = LoadResponse(status=0, capabilities=[PluginCapability.TRANSFORMER], handlers=["application/pdf"], shouldNegotiate=False)
        return res
    
    def Terminate(self, request, context):
        thread = threading.Thread(target=terminate)
        thread.start()


        return TerminateResponse(status=0)
    
    def Transform(self, request: TransformRequest, context):
        url = request.file.url
        #Get the file at URL

        req = requests.get(url)
        body = req.content
        #Read the file
        pdf = PdfReader(BytesIO(body))
        #Get the text


        root = TransformEntry(type=TransformEntryType.GROUP)
        for page in pdf.pages:
            t = page.extract_text()
            # Split text by new line, then tab, then punctuation, then space, and make a TransformEntry for each. Nest them in for statements
            text = []
            p = TransformEntry(uid=random.randint(0,2**32-1),type=TransformEntryType.GROUP, correlation=0.25)
            for line in t.split("\n"):
                g = TransformEntry(uid=random.randint(0,2**32-1),type=TransformEntryType.GROUP, correlation=0.5)
                for tab in line.split("\t"):
                    l = TransformEntry(uid=random.randint(0,2**32-1),type=TransformEntryType.GROUP, correlation=0.75)
                    for punct in tab.split("."):
                        sentence = TransformEntry(uid=random.randint(0,2**32-1),type=TransformEntryType.GROUP, correlation=1)
                        for space in punct.split(" "):
                            sentence.children.append(TransformEntry(uid=random.randint(0,2**32-1),type=TransformEntryType.STRING, contents=space, correlation=1))
                        l.children.append(sentence)
                    g.children.append(l)
                p.children.append(g)
                    
            
            root.children.append(p)
        
        rand_id = random.randint(0, 2**32-1)
        serialized_root = MessageToJson(root)
        self.entries[rand_id] = serialized_root
        # Write the contents to a file in the OS temp directory

        #file = open("{}.cache".format(rand_id), "w")
        #file.write(serialized_root)
        #print("Wrote to file",serialized_root)
        #file.close()


        print("Added entry with id",rand_id)
        print("Current entries",self.entries)
        res = TransformResponse(file=request.file, url="http://localhost:{}/dl/{}".format(self.httpPort,rand_id))



        return res
    
    
        
    def serveGrpc(self):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        add_PluginServicer_to_server(self, server)
        add_TransformerServicer_to_server(self, server)
        server.add_insecure_port('[::]:{}'.format(self.grpcPort))
        server.start()
        server.wait_for_termination()


        

def main():
    print("Starting")
    grpcPort, httpPort = sys.argv[1], sys.argv[2]
    print("GRPC Port ",grpcPort," HTTP Port ",httpPort)
    transformer = PDFTransformer()
    transformer.grpcPort = int(grpcPort)
    transformer.httpPort = int(httpPort)
    t1 = threading.Thread(target=transformer.start_http_server)
    t1.start()
    t2 = threading.Thread(target=transformer.serveGrpc)
    t2.start()
    t1.join()
    t2.join()




if __name__ == "__main__":
    main()
    

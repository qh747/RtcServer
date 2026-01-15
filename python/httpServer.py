from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import sys

class JSONRequestHandler(BaseHTTPRequestHandler):
    
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.end_headers()
    
    def do_OPTIONS(self):
        self._set_headers()
    
    def do_POST(self):
        # 获取请求长度
        content_length = int(self.headers['Content-Length'])
        
        # 读取请求体
        post_data = self.rfile.read(content_length)
        
        try:
            # 解析JSON
            json_data = json.loads(post_data.decode('utf-8'))
            print(f"Receive request: {json_data}")
            
            # 准备响应数据
            response_data = {
                "code": 0,
                "msg": "Process request success",
                "received_data": json_data  # 可选
            }
            
            # 发送响应
            self._set_headers()
            self.wfile.write(json.dumps(response_data).encode('utf-8'))
            
        except json.JSONDecodeError as e:
            error_response = {
                "code": -1,
                "msg": "Json format invalid"
            }
            self.send_response(400)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(error_response).encode('utf-8'))
        except Exception as e:
            error_response = {
                "code": -2,
                "msg": f"Server internal error: {str(e)}"
            }
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(error_response).encode('utf-8'))

def run_server(port=8093):
    server_address = ('', port)
    httpd = HTTPServer(server_address, JSONRequestHandler)
    print(f"Server start. listen port: {port}")
    print(f"Use ctrl+c to stop")
    httpd.serve_forever()

if __name__ == '__main__':
    # 获取端口参数，默认为8093
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8093
    run_server(port)

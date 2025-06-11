docker run -d --name=consul-server \
  -p 8500:8500 \
  -p 8600:8600/udp \
  consul agent -server -ui -node=my-consul -bootstrap-expect=1 -client=0.0.0.0

# 访问
# http://localhost:8500
# hack/net

This is a temporary workspace to debug and experiment with the virtual network
stack deployed to Cloudflare Containers that Apptron uses. It's pretty similar
to the `session` component, but focusing just on the network stack. It can be
connected to via websocket and used by Wanix/Apptron by providing it as an 
alternative network URL. 

Key areas of exploration:
* why is throughput so slow on Cloudflare? 
* how can we control throughput for limiting free users
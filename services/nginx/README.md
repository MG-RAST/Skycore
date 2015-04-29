

Build image:
```bash
git clone https://github.com/wgerlach/Skycore.git
cd Skycore/services/nginx/docker
docker rm -f mgrast_nginx ; docker rmi mgrast/nginx
docker build  --no-cache -t mgrast/nginx .
```

Start nginx via confd
```bash
docker run -d -p 80:80 --name mgrast_nginx mgrast/nginx
or
docker run -d -p 80:80 --name mgrast_nginx mgrast/nginx <cmd>
```
or if you want to pull from git first:
```bash
docker run -d -p 80:80 --name mgrast_nginx mgrast/nginx bash -c "cd /Skycore && git pull && /Skycore/services/nginx/confd/run_confd.sh"
```

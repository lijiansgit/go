# release
注意：
线上环境对应目录需要提前创建好
预发布:
mkdir /data/web_backup
mkdir /data/web/file_tmp
生产环境:
ssh ip 'mkdir /data/web_one'
ssh ip 'mkdir /data/web_two'
ssh ip 'ln -sv /data/web_one /data/web'
ssh ip 'chown -R  www.www /data/web/'




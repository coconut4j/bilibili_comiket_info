# bilibili_comiket_info

-d 参数是debug 打印参数及控制台输出

-mode=local是控制台输出 file是输出文件夹（默认当前文件夹）

-s 是清洗数据（只输出预售中即为开始售票的漫展）

-area=sz是输出深圳漫展的 可以改成gz

默认都是按unix时间排序的

参考参数为 ComicS.exe -d  -area=sz -mode=file -s

linux下同理
目前没写grpc部分 有服务器托管再改daemon

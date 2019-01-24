# 有些还没有完善，会输出乱码

#查看sstable文件
./lbbleveldb -index 1 -file 000033.ldb
#查看log文件
./lbbleveldb -index 2 -file 044860.log
#查看manifest文件
./lbbleveldb -index 3 -file MANIFEST-000000

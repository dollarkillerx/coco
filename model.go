package coco

type Conf struct {
	Path        string // 文档存储路径
	MaxSize     uint64 // 单文件最大存储个数  (文件越小在并行查询下越快)
	MaxReadPool uint64 // 最大读取线程数
	//MaxWritePool 采用顺序写入 理论能跑满当前磁盘性能
}

func NewDefaultConfig(path string) *Conf {
	return &Conf{
		Path:        path,
		MaxSize:     5000, // 默认单文件5000
		MaxReadPool: 300,  // 默认读取线程数300
	}
}

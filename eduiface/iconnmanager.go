package eduiface

/*
	连接管理抽象层
 */
type IConnManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除连接
	Get(connID uint32) (IConnection, error) //利用ConnID获取链接
	Len() int                               //获取当前最大连接数目
	ClearConn()															//关闭所有连接
}

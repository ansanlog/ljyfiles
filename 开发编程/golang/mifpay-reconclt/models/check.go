package main

import (
	"fmt"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_"strconv"
	"time"
	"sync"
	"log"
)

//可以改写成配置文件
var (
	dbhost string = "172.16.18.6:3306" //连接目标主机
	dbuser string = "root"           //数据库用户名
	dbpassword string = "123456"           //数据库密码
	db string = "test"            //数据库名字
)

var (
	wg sync.WaitGroup  //从mysql获取数据并发同步
	awg sync.WaitGroup  //添加到集合并发同步
	swg sync.WaitGroup  //集合操作并发同步
)

type Sbtest1 struct {
	Id  int
	K   int
	C   string
	Pad string
}

type Check_timed_task struct {
	Id             int
	Name           string
	Create_time    time.Time
	Type           int
	Separate_time  int
	Status         int
	Remark         string
	Func           string
	func_type      int
	Start_time     time.Time
	Next_exec_time time.Time
	Last_exec_time time.Time
	Is_available   int
}
type arr_check_timed_task struct {
	data []Check_timed_task
}

type Set map[Sbtest1]bool

type tblarr struct {
	data []Sbtest1
}

func init() {
	//注册mysql Driver
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//构造conn连接
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ")/" + db + "?charset=utf8"
	/*fmt.Println(conn)*/
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)
	orm.SetMaxIdleConns("default", 30)
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Sbtest1), new(Check_timed_task))
}

func main() {
	for {
		time.Sleep(10 * time.Second)
		execTask()
	}
}

func execTask() {
	/*orm.Debug=true*/
	log.Println("开始执行任务.......")
	check_timed_task := arr_check_timed_task{}
	//读取任务列表
	check_timed_task.getCheckTimedTask("select id,name,create_time,type,separate_time,status,remark,func,func_type,start_time,next_exec_time,last_exec_time,is_available from check_timed_task where (status = 0 and start_time < now()) or (status = 1 and next_exec_time < now());")
	/*fmt.Println(len(check_timed_task.data))*/

	for _, v := range check_timed_task.data {
		if v.Func == "充值渠道" {
			startTime := time.Now().UnixNano()
			sch1 := make(chan Set, 1)
			sch2 := make(chan Set, 1)
			sch3 := make(chan Set, 1)

			sbtest1, sbtest1_copy := tblarr{}, tblarr{}

			for i := 0; i < 2; i++ {
				wg.Add(1)
				if i == 1 {
					go sbtest1.getData("sbtest1")
				} else {
					go sbtest1_copy.getData("sbtest1_copy")
				}
			}

			wg.Wait()

			endTime1 := time.Now().UnixNano()
			rtt1 := int((endTime1 - startTime) / 1000000)
			fmt.Printf("data total time cost %v ms\n", rtt1)

			ssbtest1, ssbtest1_copy := make(Set, len(sbtest1.data)), make(Set, len(sbtest1_copy.data))

			for a := 0; a < 2; a++ {
				awg.Add(1)
				if a == 0 {
					go ssbtest1.Add(sbtest1.data)
				} else {
					go ssbtest1_copy.Add(sbtest1_copy.data)
				}

			}
			awg.Wait()

			endTime2 := time.Now().UnixNano()
			rtt2 := int((endTime2 - startTime) / 1000000)
			fmt.Printf("set total time cost %v ms\n", rtt2)

			for s := 0; s < 3; s++ {
				swg.Add(1)
				if s == 0 {
					go sethandle("Minus", ssbtest1, ssbtest1_copy, sch1)
				} else if s == 1 {
					go sethandle("Minus", ssbtest1_copy, ssbtest1, sch2)
				} else {
					go sethandle("Intersect", ssbtest1, ssbtest1_copy, sch3)
				}
			}
			swg.Wait()

			fmt.Println(<-sch1)
			fmt.Println(<-sch2)
			fmt.Println(len(<-sch3))

			endTime := time.Now().UnixNano()
			rtt := int((endTime - startTime) / 1000000)
			fmt.Printf("total time cost %v ms\n", rtt)

		} else {
			println("没有这个函数")
		}
	}
}

func sethandle(handletype string, tbl1, tbl2 Set, tbl chan Set) {
	if handletype == "Minus" {
		tbl <- Minus(tbl1, tbl2) // 获取差集
	} else if handletype == "Intersect" {
		tbl <- Intersect(tbl1, tbl2)  // 获取交集
	}
	swg.Done()
}

func (s *tblarr)getData(tbl string) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
			log.Printf("执行任务时发生错误:%s", err)
		}
	}();

	o := orm.NewOrm()  //注册新的orm
	o.Using("dafault") //使用数据库，默认default
	num, err := o.Raw("SELECT id,k,c,pad FROM " + tbl).QueryRows(&s.data)
	if err != nil {
		return
	} else {
		fmt.Println("user nums: ", num)
	}
	wg.Done()
}

func (arr_check_timed_task *arr_check_timed_task)getCheckTimedTask(tbl string) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
			log.Printf("执行任务时发生错误:%s", err)
		}
	}();

	o := orm.NewOrm()  //注册新的orm
	o.Using("dafault") //使用数据库，默认default
	num, err := o.Raw(tbl).QueryRows(&arr_check_timed_task.data)
	if err != nil {
		fmt.Println("arr_check_timed_task出错了", err)
		return
	} else {
		fmt.Println("user nums: ", num)
	}
}
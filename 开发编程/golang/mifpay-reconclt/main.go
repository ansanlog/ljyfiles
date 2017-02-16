package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/BurntSushi/toml"
	"sync"
	"log"

	"mifpay-reconclt/bo"
	"mifpay-reconclt/com"
	"github.com/astaxie/beego/utils"
)

const (
	CycleNUM int = 2 //获取表数据循环次数
	SetNUM int = 3 //集合操作两次差集一次交集次数
)

var (
	DBconf bo.DBconf    //配置实体
	wg sync.WaitGroup   //从mysql获取数据并发同步
	awg sync.WaitGroup  //添加到集合并发同步
	swg sync.WaitGroup  //集合操作并发同步
)

type check_timed_task_arr struct {
	data []bo.Check_timed_task
}

type tblarr struct {
	data []bo.Sbtest1
}

type arrToSet com.Set

func init() {
	if _, err := toml.DecodeFile("dbconf.toml", &DBconf); err != nil {
		log.Fatal(err)
		return
	}

	//注册mysql Driver
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//构造conn连接
	conn := DBconf.DBuser + ":" + DBconf.DBpassword + "@tcp(" + DBconf.DBhost + ")/" + DBconf.DB + "?charset=utf8"
	/*fmt.Println(conn)*/
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)
	orm.SetMaxIdleConns("default", 30)
	// 需要在init中注册定义的model
	orm.RegisterModel(new(bo.Sbtest1), new(bo.Check_timed_task))
}

func main() {
	/*startTime := time.Now().UnixNano()*/
	sch1 := make(chan com.Set, 1)
	sch2 := make(chan com.Set, 1)
	sch3 := make(chan com.Set, 1)

	//并发操作：获取对账数据
	sbtest1, sbtest1_copy := tblarr{}, tblarr{}

	for i := 0; i < CycleNUM; i++ {
		wg.Add(1)

		switch i {
		case 0:
			go sbtest1.getCheckDataFromMySQL("SELECT id,k,c,pad FROM sbtest1")
		case 1:
			go sbtest1_copy.getCheckDataFromMySQL("SELECT id,k,c,pad FROM sbtest1_copy")
		default:
			fmt.Println("超出范围啦！！")
		}
	}
	wg.Wait()

	//数组转化成集合
	ssbtest1, ssbtest1_copy := make(com.Set, len(sbtest1.data)), make(com.Set, len(sbtest1_copy.data))
	for i:=0;i< CycleNUM; i++ {


	}
	ssbtest1.Add(sbtest1.data)
	ssbtest1_copy.Add(sbtest1_copy.data)
	//并发操作：两次差集一次交集
	for s := 0; s < SetNUM; s++ {
		swg.Add(1)

		switch s {
		case 0:
			go setHandle("Minus", ssbtest1, ssbtest1_copy, sch1)
		case 1:
			go setHandle("Minus", ssbtest1_copy, ssbtest1, sch2)
		case 2:
			go setHandle("Intersect", ssbtest1, ssbtest1_copy, sch3)
		default:
			fmt.Println("超出范围啦！！")
		}
	}
	swg.Wait()

	fmt.Println(<-sch1)
	fmt.Println(<-sch2)
	fmt.Println(len(<-sch3))
}

func (s arrToSet)setArrToSet(arr *tblarr)  {
	newSet := make(com.Set, len(arr.data))
	newSet.Add(arr.data)
	s=newSet
}

func setHandle(handleType string, tbl1, tbl2 com.Set, tbl chan com.Set) {

	switch handleType {
	case "Minus":
		tbl <- utils.Minus(tbl1, tbl2) // 获取差集
	case "Intersect":
		tbl <- utils.Intersect(tbl1, tbl2)  // 获取交集
	default:
		fmt.Println("没有这项操作！！")
	}

	swg.Done()
}

func (check_timed_task_arr *check_timed_task_arr)getCheckTimedTask(tbl string) {
	//异常处理
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
			log.Printf("执行获取定时任务数据时发生错误:%s", err)
		}
	}();

	o := orm.NewOrm()  //注册新的orm
	o.Using("dafault") //使用数据库，默认default
	_, err := o.Raw(tbl).QueryRows(&check_timed_task_arr.data)
	if err != nil {
		fmt.Println("执行getCheckTimedTask时出现异常了", err)
		return
	}
}

func (s *tblarr)getCheckDataFromMySQL(strSQL string) {
	//异常处理
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
			log.Printf("执行获取定时任务数据时发生错误:%s", err)
		}
	}();

	o := orm.NewOrm()  //注册新的orm
	o.Using("dafault") //使用数据库，默认default
	_, err := o.Raw(strSQL).QueryRows(&s.data)
	if err != nil {
		return
	}
	wg.Done()
}
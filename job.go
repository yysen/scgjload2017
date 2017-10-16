package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"fmt"

	"os"
	"path/filepath"

	"archive/zip"

	"bufio"

	"github.com/linlexing/dbx/data"
	"github.com/linlexing/dbx/ddb"
	"github.com/linlexing/dbx/schema"
	"github.com/pborman/uuid"
	"github.com/robfig/cron"
)

const batchNum = 500

var (
	jobs    = cron.New()
	running = false
	jobRun  = &sync.Mutex{}
)

func taskRun() {
	if running {
		return
	}
	running = true
	defer func() {
		running = false
	}()
	dlog.Println("before lock")
	jobRun.Lock()
	defer jobRun.Unlock()
	dlog.Println("start job")
	err := buildDataFile()
	if err != nil {
		dlog.Error(err)
		return
	}
	//然后开始上传
	uploadAll()

	dlog.Println("job finished")
}

//创建zip文件，并将templa文件夹中的文件复制到zip文件中
func createNewZipFile() (*os.File, *zip.Writer, error) {
	//添加在workDir后添加out路径
	outPath := filepath.Join(workDir, "out")
	//使用指定的权限和名称创建一个目录，包括任何必要的上级目录，并返回nil，否则返回错误
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return nil, nil, err
	}
	//确定文件名过程为：
	//out目录中没有同名文件
	//upload目录中也没有同名文件
	var fileName string
	for i := 1; ; i++ {
		fileName = fmt.Sprintf("gsdata_%s_%s_%06d.zip", time.Now().Format("20060102"),
			vconfig.AreaCode, i)
		var not1, not2 bool
		//Stat返回一个描述name指定的文件对象的FileInfo。
		//如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接指向的文件的信息，
		//本函数会尝试跳转该链接
		//IsNotExist返回一个布尔值说明该错误是否表示一个文件或目录不存在
		if _, err := os.Stat(filepath.Join(workDir, vconfig.FinishOut, fileName)); os.IsNotExist(err) {
			not1 = true
		} else if err != nil {
			return nil, nil, err
		}
		if _, err := os.Stat(filepath.Join(workDir, "out", fileName)); os.IsNotExist(err) {
			not2 = true
		} else if err != nil {
			return nil, nil, err
		}
		if not1 && not2 {
			break
		}
	}
	//创建文件名为fileName的zip
	file, err := os.Create(filepath.Join(outPath, fileName))
	if err != nil {
		return nil, nil, err
	}
	//得到一个将zip文件写入file的*Writer
	zipw := zip.NewWriter(file)
	//先复制模板文件
	//返回template指定的目录的目录信息的有序列
	files, err := ioutil.ReadDir(filepath.Join(workDir, "template"))
	if err != nil {
		return nil, nil, err
	}
	//将template文件夹里的文件复制到Zip中
	for _, f := range files {
		//使用给出的文件名添加一个文件进zip文件。本方法返回的w是一个io.Writer接口（用于写入新添加文件的内容）
		w, err := zipw.Create(f.Name())
		if err != nil {
			return nil, nil, err
		}
		//ReadFile 从filename指定的文件中读取数据并返回文件的内容
		bys, err := ioutil.ReadFile(filepath.Join(workDir, "template", f.Name()))
		if err != nil {
			return nil, nil, err
		}
		//通过w向文件中写入bys
		if _, err = w.Write(bys); err != nil {
			return nil, nil, err
		}
	}
	//w, err := zipw.Create("ent_info.dat")  bufio.NewWriter(w),

	//bufio.NewWriter创建一个具有默认大小缓冲、写入w的*Writer。
	return file, zipw, err
}

//打开主表和对应对的影子表
func openDB(fieldSize []int, tableName, shadowTableName, ID string) (ddb.DB, *data.Table, *data.Table, error) {
	db, err := ddb.Openx(vconfig.Driver, vconfig.DBURL)
	if err != nil {
		return nil, nil, nil, err
	}
	tab, err := data.OpenTable(db.DriverName(), db, tableName)
	if err != nil {
		return nil, nil, nil, err
	}
	//必须全部是string类型
	for _, col := range tab.Columns {
		if col.Type != schema.TypeString {
			return nil, nil, nil, fmt.Errorf("column %s type not is string", col.Name)
		}
	}
	//配置文件里的字段没有连接字段，在此加上
	size := len(fieldSize) + 1
	if size != len(tab.Columns) {
		return nil, nil, nil, fmt.Errorf("field size %d <> column length %d", size, len(tab.Columns))
	}
	tab.Name = shadowTableName
	tab.PrimaryKeys = []string{strings.ToUpper(ID)}
	//自动更新影子表的结构
	if err = tab.Table.Update(db.DriverName(), db); err != nil {
		return nil, nil, nil, err
	}
	tab.Name = tableName
	shadowTable, err := data.OpenTable(db.DriverName(), db, shadowTableName)
	if err != nil {
		return nil, nil, nil, err
	}
	return db, tab, shadowTable, nil
}

//打开子表，并根据关联条件查找记录，返回一个带有主表uuid的切片
func openSonTable(db ddb.DB, uuid interface{}, idName, id, tableName string) ([][]interface{}, error) {
	//查找子表记录内容
	rows, err := db.Query(fmt.Sprintf("select * from %s where %s='%s'", tableName, idName, id))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	sonRow := [][]interface{}{}
	//得到所有字段名
	sonColumn, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colCount := len(sonColumn)
	if rows.Next() {
		row := make([]interface{}, colCount)
		for i := range row {
			row[i] = new(interface{})
		}
		if err = rows.Scan(row...); err != nil {
			return nil, err
		}
		//将interface类型转换为string
		line := make([]interface{}, colCount)
		for i := range row {
			line[i] = *(row[i].(*interface{}))
		}
		//将主表uuid放进查询结果中
		line[2] = uuid
		sonRow = append(sonRow, line)
	}
	return sonRow, nil
}

//searchTable 遍历表，找出与影子表有差异的数据，并将差异数据更新到影子表
func searchTable(maxrows int, ID string, db ddb.DB, tab, shadowTable *data.Table, cb func(icount int,
	diffrows [][]interface{}) error) error {
	saveToShadowTable := func(diffRows [][]interface{}) error {
		//保存到影子表中
		for _, line := range diffRows {
			row := map[string]interface{}{}
			for i, col := range shadowTable.Columns {
				row[col.Name] = line[i]
			}
			if err := shadowTable.Save(row); err != nil {
				return err
			}
		}
		return nil
	}
	rows, err := db.Query(fmt.Sprintf("select %s from %s", ID, tab.Name))
	if err != nil {
		return err
	}
	defer rows.Close()
	//差异比较的数据量
	icount := 1
	//上传数据的数量
	rowsNum := 0
	pks := []interface{}{}
	for rows.Next() {
		var pk interface{}
		//得到每一行的第一个，即主键的值
		if err = rows.Scan(&pk); err != nil {
			return err
		}
		pks = append(pks, pk)
		if icount%batchNum == 0 {
			diffRows, err := queryDiff(db, tab, shadowTable.Name, pks)
			if err != nil {
				return err
			}
			rowsNum += len(diffRows)
			pks = nil
			if err := cb(rowsNum, diffRows); err != nil {
				return err
			}
			if err := saveToShadowTable(diffRows); err != nil {
				return err
			}

		}
		icount++
		//限制每次上传数据的数量
		if rowsNum >= maxrows {
			break
		}
	}
	if len(pks) > 0 {
		diffRows, err := queryDiff(db, tab, shadowTable.Name, pks)
		if err != nil {
			return err
		}
		rowsNum += len(diffRows)
		pks = nil
		if err := cb(rowsNum, diffRows); err != nil {
			return err
		}
		if err := saveToShadowTable(diffRows); err != nil {
			return err
		}
	}
	return nil
}

func writeLine(w *bufio.Writer, FieldSize []int, data []interface{}) error {
	for i, v := range data {
		var str string
		switch tv := v.(type) {
		case string:
			str = tv
		case []byte:
			str = string(tv)
		case nil:
		default:
			return fmt.Errorf("%T not in string", v)
		}
		str = strings.Replace(
			strings.Replace(str, "\r", " ", -1),
			"\n", " ", -1)
		//得到str的字符长度，而不是字节
		rstr := []rune(str)
		if len(rstr) < FieldSize[i] {
			str = str + strings.Repeat(" ", FieldSize[i]-len(rstr))
		} else if len(rstr) > FieldSize[i] {
			str = string(rstr[:FieldSize[i]])
		}
		if _, err := w.WriteString(str); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}

//向文件中写入数据
func buildDataFile() error {
	file, zipw, err := createNewZipFile()
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	defer zipw.Close()

	//打开主表，并对数据进行简单处理
	tables := vconfig.Table
	for _, maintable := range tables {
		//创建dat文件，并得到一个向该文件写的流
		datFile, err := zipw.Create(fmt.Sprint(maintable.Name, ".dat"))
		if err != nil {
			log.Println(err)
			return err
		}
		//得到一个向dat文件写的缓冲流
		datw := bufio.NewWriter(datFile)
		defer datw.Flush()
		/*打开对应的主表*/
		db, tab, shadowTable, err := openDB(maintable.FieldSize, maintable.Name, maintable.Shadowtable, maintable.ID)
		if err != nil {
			log.Println(err)
			return err
		}
		dateTimeFields := []int{}
		for i, col := range tab.Columns {
			switch col.Name {
			case "数据修改时间", "数据上传时间":
				dateTimeFields = append(dateTimeFields, i)
			}
		}
		icount := 0
		if err := searchTable(maintable.Maxrows, maintable.ID, db, tab, shadowTable, func(i int, rows [][]interface{}) error {
			if len(rows) > 0 {
				icount += len(rows)
				dlog.Println("rownum:", i, "write", len(rows), "rows,total:", icount)
			}
			for _, line := range rows {

				//不复制，会影响回写影子表,且去除元数据中的主键
				newLine := make([]interface{}, len(line)-1)
				copy(newLine, line[1:])
				//设置主表uuid的值
				newLine[1] = interface{}(hex.EncodeToString(uuid.NewUUID()))
				//设置时间字段
				for _, idx := range dateTimeFields {
					newLine[idx-1] = time.Now().Format("20060102150405")
				}
				if err := writeLine(datw, maintable.FieldSize, newLine); err != nil {
					log.Println(err)
					return err
				}
				id, isString := line[0].(string)
				//判断主键是否转换为string类型
				if !isString {
					continue
				}
				//判断有无子表,
				if len(maintable.Details) > 0 {
					for _, sonTable := range maintable.Details {
						sonRow, err := openSonTable(db, newLine[1], sonTable.PreID, id, sonTable.Name)
						if err != nil {
							log.Println(err)
							return err
						}
						//判断子表是否有相应内容
						if len(sonRow) == 0 {
							continue
						}
						for _, sonLine := range sonRow {
							newSonLine := make([]interface{}, len(sonLine))
							newSonLine = sonLine[1:]
							//写入子表内容
							if err := writeLine(datw, sonTable.Sizes, newSonLine); err != nil {
								log.Println(err)
								return err
							}
						}
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
		//创建审核结果文件和审核状态文件
		_, err = zipw.Create(fmt.Sprint(maintable.Result, ".dat"))
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = zipw.Create(fmt.Sprint(maintable.Status, ".dat"))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

//queryDiff 返回一个新增、变更的记录内容
func queryDiff(db ddb.DB, table *data.Table, shadowtable string,
	pkvalues []interface{}) ([][]interface{}, error) {
	str := fmt.Sprintf("%s in(?)", table.PrimaryKeys[0])
	where, params, err := data.In(str, pkvalues)
	if err != nil {
		return nil, err
	}
	//两个表的where
	params = append(params, params...)
	strSQL := data.Find(db.DriverName()).Minus(table.FullName(), where,
		shadowtable, where, table.PrimaryKeys, table.ColumnNames)
	//log.Println("strsql ", strSQL)
	rows, err := db.Query(strSQL, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rev := [][]interface{}{}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colCount := len(cols)
	for rows.Next() {
		row := make([]interface{}, colCount)
		for i := range row {
			row[i] = new(interface{})
		}
		if err = rows.Scan(row...); err != nil {
			return nil, err
		}
		line := make([]interface{}, colCount)
		for i := range row {
			line[i] = *(row[i].(*interface{}))
		}

		rev = append(rev, line)
	}
	return rev, nil
}

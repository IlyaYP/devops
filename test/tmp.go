package test

import (
	"context"
	"fmt"
	"github.com/IlyaYP/devops/storage/inmemory"
)

func testStore(st *inmemory.Storage) {
	println(st)

	if err := st.PutMetric(context.Background(), "gauge", "aaa", "111.111"); err != nil {
		println(err.Error())
		return
	}

	if err := st.PutMetric(context.Background(), "counter", "bbb", "222"); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := st.PutMetric(context.Background(), "gauge", "ccc", "333.333"); err != nil {
		fmt.Println(err.Error())
		return
	}

	if v, err := st.GetMetric(context.Background(), "gauge", "ccc"); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println(v)
	}
	ret := st.ReadMetrics(context.Background())
	fmt.Println(ret)

}

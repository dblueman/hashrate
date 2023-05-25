package hashrate


import (
   "fmt"
   "testing"
   "time"
)

func Test(t *testing.T) {
   tb := New[string](4., 6.)

   for i := 0; i < 4; i++ {
      v := tb.Outlier("first")
      fmt.Printf("%v %+v\n", v, tb.events["first"])
   }

   time.Sleep(time.Second)
   fmt.Printf("%+v\n", tb.events["first"])

   for i := 0; i < 10; i++ {
      v := tb.Outlier("first")
      fmt.Printf("%v %+v\n", v, tb.events["first"])

      if v {
         t.Error("expected no outlier")
      }

      time.Sleep(time.Millisecond * time.Duration(150))
   }

   var v bool

   for i := 0; i < 10; i++ {
      v = tb.Outlier("second") // overwrite state
      fmt.Printf("%v %+v\n", v, tb.events["second"])
      time.Sleep(time.Millisecond * time.Duration(100))
   }

   if !v {
      t.Error("expected outlier")
   }
}

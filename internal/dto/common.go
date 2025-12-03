// internal/dto/common.go
package dto

// c.ShouldBindUri()			uri:"age"
// c.ShouldBindQuery()		form:"age"
// c.ShouldBindJSON()			json:"age"

// required				 // 必填
// omitempty       // 可选（默认就是）

// min=6           // 字符串/数字最小
// max=32          // 最大
// len=11          // 长度必须等于

// email           // 邮箱格式
// alpha           // 只能是字母
// alphanum        // 只能是字母和数字
// numeric         // 只能是数字

// contains=@      // 必须包含 @
// oneof=1 2 3     // 只能是这些值之一
// gt=0            // 大于 0
// gte=18          // 大于等于 18

// oneof=1 2 3     // 只能是这些值之一

// eqfield=Password	 // 等于 Password
// gtfield=StartTime // 比 StartTime 大

<script setup>
import {DeleteRegistry, SearchRegistry} from "../../wailsjs/go/main/App";
import {EventsOn} from "../../wailsjs/runtime";
import {computed, h, nextTick, onMounted, reactive, ref} from 'vue';
import {NTag, useDialog, useMessage} from "naive-ui";
import {da} from "date-fns/locale";

const softwareName = ref('')
const data = ref([]);
const percentage = ref(0)
// const data = ref([
//   {id:0,tags:["medium"],key:"moba",path:"/13/123"},
//   {id:1,tags:["medium"],key:"moba",path:"/13/123"},
//   {id:2,tags:["medium"],key:"xmanager",path:"/13/123"},
//   {id:3,tags:["high"],key:"moba",path:"/13/123"},
//   {id:4,tags:["low"],key:"xshell",path:"/13/123"},
// ])
//generateDirtyData(1000)
function generateDirtyData(count) {
  const tags = ["low", "medium", "high"]; // 预定义的标签数组
  const keys = ["moba", "xmanager", "xshell"]; // 预定义的键数组
  const path = "/13/123"; // 固定的路径
  for (let i = 0; i < count; i++) {
    const randomTags = tags[Math.floor(Math.random() * tags.length)]; // 随机选择一个标签
    const randomKey = keys[Math.floor(Math.random() * keys.length)]; // 随机选择一个键
    data.value.push({
      id: i, // 使用循环索引作为id
      tags: [randomTags], // 将随机标签放入数组
      key: randomKey, // 使用随机键
      path: path // 使用固定的路径
    });
  }
}

const input = ref('');
//const selectedRows = ref([]);

const selectedRowsRef = ref([]);

const rowProps = (row) => {
  return {
    onContextmenu: (e) => {
      e.preventDefault();
      showDropdownRef.value = false;
      nextTick().then(() => {
        showDropdownRef.value = true;
        xRef.value = e.clientX;
        yRef.value = e.clientY;
        //selectedRowsRef.value = []
        //selectedRowIndex.value = index;
        //console.log(selectedRows.value)
      });
    },
    style: {
      cursor: 'pointer',
    },
  };
};

const handleCheckKeys = (rowKeys) =>{
  // 清空数组
  selectedRowsRef.value = []
  selectedRowsRef.value = rowKeys
  //rowKeys=rowKeys.filter(rowKey => rowKey !== selectedRowsRef.value[0]);
  // rowKeys.forEach(rowKey => {
  //   selectedRowsRef.value.push(rowKey);
  // });
}

const rowKey = (row) => row.id

const columns = ref([
  {
    type: "selection",
  },
  {
    title: 'ID',
    key: 'id',
    width: 50,
    sorter: 'default'
  },
  {
    title: 'level',
    key: 'tags',
    width: 80,
    render(row) {
      const tagKey = row.tags[0]; // 获取 tags 数组中的第一个元素
      let type = "info"; // 默认type为info  
        
        // 根据tagKey的值设置不同的type  
        if (tagKey === 'high') {  
          type = "error"; 
        } else if (tagKey === 'medium') {  
          type = "warning";  
        } else if (tagKey === 'low') {  
          type = "info";   
        }  
      return h(
          NTag,
          {
            style: {
              marginRight: "6px"
            },
            type: type,
            bordered: false
          },
          {
            default: () => tagKey
          }
      );
    },
    sorter: (a, b) => {
      // 获取 tags 数组中的第一个元素进行排序
      const tagTextA = a.tags[0];
      const tagTextB = b.tags[0];
      return tagTextA.localeCompare(tagTextB);
    }
  },
  {
    title: '软件名称',
    key: 'key',
    width: 200,
    resizable: true,
    sorter: 'default'
  },
  {
    title: '注册表路径',
    key: 'path',
    resizable: true,
    sorter: 'default'
  },
])

const play = (row)=> {
  message.info(`Play ${row.title}`);
}


const options = [
  // {
  //   label: "编辑",
  //   key: "edit"
  // },
  {
    label: () => h("span", { style: { color: "red" } }, "选择的"+ selectedRowsRef.value.length +"个注册表项删除"),
    key: "delete"
  }
];

const xRef = ref(0);
const yRef = ref(0);
const showDropdownRef = ref(false);
const x=  xRef;
const y=  yRef;

// 在选中行的索引上设置一个ref
//const selectedRowIndex = ref(-1);

const dialog = useDialog()
// let sure = false
// const handleConfirm = () => {
//   dialog.warning({
//     title: "警告",
//     content: "确定删除？",
//     positiveText: "确定",
//     negativeText: "不确定",
//     onPositiveClick: () => {
//       sure = true
//     },
//     onNegativeClick: () => {
//       sure = false
//     }
//   });
// }

const deleteNumber = ref(0)

const deleteSelect = ()=>{
  if(data.value.length===0){
    message.info("无可删除项！")
    return
  }
  if(selectedRowsRef.value.length===0){
    message.info("请先选中行！")
    return;
  }
  showDropdownRef.value = false;
  dialog.warning({
    title: "警告",
    content: "确定删除？",
    positiveText: "确定",
    negativeText: "不确定",
    onPositiveClick: async () => {
      for (const row of data.value) {
        if (selectedRowsRef.value.includes(row.id)) {
          const result = await DeleteRegistry(row.path);
          if (result === "success") {
            deleteNumber.value++
            data.value = data.value.filter(item => item.id !== row.id);
          } else {
            dialog.error({
              title: "删除失败，请手动去注册表删除该项",
              content: result,
              positiveText: "确定",
              onPositiveClick: () => {
              }
            })
            console.log(result);
          }
        }
      }
      selectedRowsRef.value = [];
    },
    onNegativeClick: () => {

    }
  });
}

const handleSelect = (key)=> {
  showDropdownRef.value = false;
  if (key === 'delete') {
    //调用后端go 删除该注册表项
      // 从数据中移除选中的行
    dialog.warning({
      title: "警告",
      content: "确定删除？",
      positiveText: "确定",
      negativeText: "不确定",
      onPositiveClick: async () => {
        for (const row of data.value) {
          if (selectedRowsRef.value.includes(row.id)) {
            const result = await DeleteRegistry(row.path);
            if (result === "success") {
              deleteNumber.value++
              data.value = data.value.filter(item => item.id !== row.id);
            } else {
                dialog.error({
                  title: "删除失败，请手动去注册表删除该项",
                  content: result,
                  positiveText: "确定",
                  onPositiveClick: () => {
                  }
                })
              console.log(result);
            }
          }
        }
        selectedRowsRef.value = [];
      },
      onNegativeClick: () => {

      }
    });
  }
}

const onClickoutside = () =>{
  showDropdownRef.value = false;
}
const  showDropdown = showDropdownRef;

const message = useMessage()

const initRegistryMap = async() =>{
  //每次点击搜索都重置
  data.value = []
  index.value = 1
  percentage.value = 0.0
  disable.value = true

  //message.info("正在扫描注册表，请稍候~");
  if (input.value.length === 0||input.value===' '||input.value==="\n") {
    //input.value = ""
    message.info("请按格式输入软件名称！")
    disable.value = false
    return
  }
  message.info("注册表初始化扫描中请稍后~")
  console.log(percentage.value)

  await SearchRegistry(input.value)

  console.log(data.value)
  disable.value = false
}

const index = ref(1)

EventsOn("SearchRegistry", e => {
  //去重
  const isPathExist =  data.value.some(item => item.path === e.path)
  let tagsArray = [e.accuracy];
  if (!isPathExist) {
    data.value.push({
      id: index.value++,
      tags: tagsArray,
      key: e.key,
      path: e.path
    });
  }

});

EventsOn("percentage", e=>{
  if(percentage.value < 100){
    percentage.value = parseInt(e)
  }

})


// EventsOn("DeleteError", e=>{
//   dialog.error({
//     title: "删除失败，请手动去注册表删除该项",
//     content: e,
//     positiveText: "确定",
//     onPositiveClick: () => {
//     }
//   })
// })


const handleChange = (event) => {
  //console.log(event)
  const files = event.target.files;
  if (files.length > 0) {
    //console.log(event)
    const file = files[0];
    readFileContent(file);
  }
};
const readFileContent = async (file) => {
  try {
    const content = await readFile(file);
    input.value = content;
  } catch (error) {
    console.error("Error reading the file:", error);
  }
};

const readFile = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const content = reader.result.replace(/\r\n/g, '\n');
      resolve(content);
    };
    reader.onerror = () => {
      reject(reader.error);
    };
    reader.readAsText(file,'UTF-8');
  });
};

const paginationReactive = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 50, 100],
  onChange: (page) => {
    paginationReactive.page = page;
  },
  onUpdatePageSize: (pageSize) => {
    paginationReactive.pageSize = pageSize;
    paginationReactive.page = 1;
  }
});

const filteredItems = computed(() => {
  return data.value.filter(item =>
      item.path.toLowerCase().includes(softwareName.value.toLowerCase())
  );
});

const disable = ref(false)

</script>

<template>
  <div class="box">
    <div class="search">
    <n-input v-model:value="input" type="textarea"
             placeholder="输入软件名称或导入黑名单
格式:
huorong
bilibili
wechat"
             :autosize="{
                      minRows: 5,
                      maxRows: 20
                    }"
             class="input-l"/>
    <n-button type="primary" class="button" @click="initRegistryMap" :disabled=disable>
      点击搜索
    </n-button>
    <div class="custom-file-upload">
      <label for="file-upload" class="custom-file-upload-label">
        <span style="font-family: inherit;">导入黑名单</span>
      </label>
      <input type="file" id="file-upload" @change="handleChange">
    </div>
    <n-progress
        v-if="percentage < 100"
        class="Pro"
        type="circle"
        :percentage=percentage
        :indicator-placement="'inside'"

    />
    </div>
    <div class="line">
    <n-row class="row">
      <n-col :span="12">
        <n-statistic label="统计数据(已删除/残留项)" :value=deleteNumber>
          <template #prefix>
            <n-icon>
              <md-save />
            </n-icon>
          </template>
          <template #suffix>
            /{{data.length}}
          </template>
        </n-statistic>
      </n-col>
    </n-row>
      <div class="delete">
      <n-button type="primary" class="button" @click="deleteSelect">
        删除选中项
      </n-button>
      </div>
    </div>
    <n-input v-model:value="softwareName" placeholder="注册表路径过滤搜索"></n-input>
    <n-data-table
        :pagination="paginationReactive"
        :checked-row-keys = "selectedRowsRef"
        @update:checked-row-keys="handleCheckKeys"
        :row-key="rowKey"
        :row-props="rowProps"
        size="small"
        :bordered="false"
        :single-line="false"
        :columns="columns"
        :data="filteredItems"
        :max-height="300"
        style="margin-top: 10px"
    >
    </n-data-table>
    <n-dropdown
        placement="bottom-start"
        trigger="manual"
        :x="x"
        :y="y"
        :options="options"
        :show="showDropdown"
        :on-clickoutside="onClickoutside"
        @select="handleSelect"
    />
  </div>
</template>

<style scoped>
.box {
  height: 100%;
  width: 90vw;

  .input-l {
    float: left;
    text-align: left;
    width: 50%;
  }

  .button {
    margin-left: 10px;
    float: left;
  }
  .Pro{
    margin: 20px;
  }
  .input {
    font-weight: bold;
    margin-top: 20px;
    float: left;
    text-align: left;
    width: 100%;
  }

  .custom-file-upload {
    padding-left: 20px;
    float: left;
    left: 20px;
    display: inline-block;
  }
}

.custom-file-upload-label {
  border-radius: 5px;
  display: inline-block;
  background-color: #18a058; /* 半透明背景颜色 */
  color: #fff;
  padding: 5px 10px;
  border: 1px solid #ccc; /* 灰色边框 */
  width: 70px;
  height: 25px;
  text-align: center;
  cursor: pointer;
}

/* 隐藏原生文件输入框 */
#file-upload {
  display: none;
}
.search{
  margin-bottom: 10px;
  display: flex;
}
.row{
  text-align: center;
}
.line{
  display: flex;
  .delete{
    text-align: center;
    margin-top: 10px;
    margin-right: 350px;
  }
}
</style>

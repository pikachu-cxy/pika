<script setup>
import {DeleteRegistry, SearchRegistry} from "../../wailsjs/go/main/App";
import {EventsOn} from "../../wailsjs/runtime";
import {computed, h, nextTick, onMounted, ref} from 'vue';
import {NTag} from "naive-ui";


const data = ref([]);
const percentage = ref(0)
// const data = ref([
//   {id:0,tags:["medium"],key:"moba",path:"/13/123"},
//   {id:1,tags:["medium"],key:"moba",path:"/13/123"},
//   {id:2,tags:["medium"],key:"xmanager",path:"/13/123"},
//   {id:3,tags:["high"],key:"moba",path:"/13/123"},
//   {id:4,tags:["low"],key:"xshell",path:"/13/123"},
// ])

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
      return h(
          NTag,
          {
            style: {
              marginRight: "6px"
            },
            type: "info",
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
  }
])

const options = [
  // {
  //   label: "编辑",
  //   key: "edit"
  // },
  {
    label: () => h("span", { style: { color: "red" } }, "选择的"+ selectedRowsRef.value.length +"行删除"),
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

const handleSelect = (key)=> {
  showDropdownRef.value = false;
  if (key === 'delete') {
    //调用后端go 删除该注册表项

    // 从数据中移除选中的行
    data.value.forEach(
        row =>{
          if (selectedRowsRef.value.includes(row.id)) {
            // 从数据中移除选中的行
            const index = data.value.findIndex(item => item.id === row.id);
            if (index !== -1) {
              DeleteRegistry(row.path);
              data.value.splice(index, 1);
            }
            // 调用后端 Go 函数来删除注册表项
            //应该在前端删除表项前，后端删除，如返回错误 则回复删除失败
            //DeleteRegistry(row.path);
          }
        });
    //data.value = data.value.filter(row => !selectedRowsRef.value.includes(row.id));
    //DeleteRegistry(row.path)
    selectedRowsRef.value = []
  }
}

const onClickoutside = () =>{
  showDropdownRef.value = false;
}
const  showDropdown = showDropdownRef;

const initRegistryMap = async() =>{
  //每次点击搜索都重置
  data.value = []
  index.value = 1
  percentage.value = 0

  //message.info("正在扫描注册表，请稍候~");
  if (input.value.length === 0) {
    //input.value = ""
    return
  }
  await SearchRegistry(input.value)

  //console.log(data.value)
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

</script>

<template>
  <n-input v-model:value="input" type="textarea"
           placeholder="输入关键字"
           :autosize="{
                      minRows: 5,
                      maxRows: 20
                    }"
           class="input-l"/>
  <n-button type="primary" class="button" @click="initRegistryMap">
    点击搜索
  </n-button>
  <n-progress
      type="line"
      :percentage=percentage
      :indicator-placement="'inside'"
  />
  <n-data-table
      :checked-row-keys = "selectedRowsRef"
      @update:checked-row-keys="handleCheckKeys"
      :row-key="rowKey"
      :row-props="rowProps"
      size="small"
      :bordered="false"
      :single-line="false"
      :columns="columns"
      :data="data"
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
</template>

<style scoped>
.selected {
  background-color: gold;
  color: #18a058;
}
:deep(.selected){
  background-color: gold;
  color: #18a058;
  cursor: pointer;
}
</style>
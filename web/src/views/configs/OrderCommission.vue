<template>
  <div class="orderCommission">
    <div v-if="!editMode">
      <ConfigTable :category="category" name="订单佣金配置" />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </div>
    <ConfigEditor
      name="添加/更新订单佣金配置"
      summary="配置各分组佣金"
      :category="category"
      :defaultValue="defaultValue"
      v-else
    >
      <template v-slot:data="configProps">
        <OrderCommissionData
          :data="configProps.form.data"
          @change="configProps.form.data = $event"
        />
      </template>
    </ConfigEditor>
  </div>
</template>

<script>
import { ORDER_COMMISSION } from "@/constants/config";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import ConfigEditor from "@/components/configs/Editor.vue";
import ConfigTable from "@/components/configs/Table.vue";
import OrderCommissionData from "@/components/configs/OrderCommissionData.vue";
export default {
  name: "OrderCommission",
  components: {
    ConfigEditor,
    ConfigTable,
    OrderCommissionData
  },
  data() {
    return {
      defaultValue: {
        category: ORDER_COMMISSION
      },
      category: ORDER_COMMISSION
    };
  },
  computed: {
    editMode() {
      return this.$route.query.mode === CONFIG_EDITE_MODE;
    }
  },
  methods: {
    add() {
      this.$router.push({
        query: {
          mode: CONFIG_EDITE_MODE
        }
      });
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.add
  margin: $mainMargin
.addBtn
  width: 100%
</style>

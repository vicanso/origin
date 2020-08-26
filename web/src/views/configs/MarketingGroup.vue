<template>
  <div class="marketingGroup">
    <div v-if="!editMode">
      <ConfigTable :category="category" name="销售分组" />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </div>
    <ConfigEditor
      name="添加/更新销售分组"
      summary="配置销售分组"
      :category="category"
      :defaultValue="defaultValue"
      v-else
    >
      <template v-slot:data="configProps">
        <MarketingGroupData
          :data="configProps.form.data"
          @change="configProps.form.data = $event"
        />
      </template>
    </ConfigEditor>
  </div>
</template>
<script>
import { MARKETING_GROUP } from "@/constants/config";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import ConfigEditor from "@/components/configs/Editor.vue";
import ConfigTable from "@/components/configs/Table.vue";
import MarketingGroupData from "@/components/configs/MarketingGroupData.vue";

export default {
  name: "MarketingGroup",
  components: {
    ConfigEditor,
    ConfigTable,
    MarketingGroupData
  },
  data() {
    return {
      defaultValue: {
        category: MARKETING_GROUP
      },
      category: MARKETING_GROUP
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

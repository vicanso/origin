<template>
  <el-cascader
    class="regionSelect"
    v-if="inited"
    clearable
    :props="props"
    @change="handleChange"
    v-model="region"
    :value="region"
    :show-all-levels="$props.showAllLevels"
  ></el-cascader>
</template>
<script>
import { mapActions, mapState } from "vuex";

const countryLevel = 0;
const chinaCode = "CN";

export default {
  name: "RegionSelect",
  props: {
    maxLevel: {
      type: Number,
      default: 4
    },
    startLevel: {
      type: Number,
      default: 0
    },
    showAllLevels: {
      type: Boolean,
      default: false
    },
    value: String
  },
  data() {
    const region = [];
    // 如果代码全是数字，直接转换为中国地址
    const { value } = this.$props;
    if (value) {
      if (/\d/.test(value)) {
        // 440111009
        // 地址代码分割：省两位，市两位，区两位，后面的为街道
        const offsets = [2, 2, 2, 3];
        region.push("CN");
        let offset = 0;
        offsets.forEach(item => {
          offset += item;
          if (value.length >= offset) {
            region.push(value.substring(0, offset));
          }
        });
      } else {
        region.push(value);
      }
    }
    return {
      inited: false,
      region,
      query: {
        limit: 100,
        status: 1,
        fields: "code,name",
        order: "-priority"
      },
      props: {
        lazy: true,
        lazyLoad: this.lazyLoad
      }
    };
  },
  computed: mapState({
    categories: state => state.region.categories || []
  }),
  methods: {
    ...mapActions(["listRegion", "listRegionCategory"]),
    handleChange(value) {
      const v = value[value.length - 1];
      this.$emit("input", v);
    },
    async lazyLoad(node, resolve) {
      const { maxLevel, startLevel } = this.$props;
      let level = node.level + startLevel;
      let category = this.categories[level].value;
      const leaf = level === maxLevel - 1;
      let parent = "";
      if (node.data) {
        parent = node.data.value;
      }
      try {
        const { regions } = await this.listRegion({
          params: Object.assign(
            {
              parent,
              category
            },
            this.query
          ),
          category
        });
        const nodes = regions.map(item => {
          let isLeaf = leaf;
          if (level === countryLevel && item.code !== chinaCode) {
            isLeaf = true;
          }
          return {
            value: item.code,
            label: item.name,
            leaf: isLeaf
          };
        });
        resolve(nodes);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    try {
      await this.listRegionCategory();
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.inited = true;
    }
  }
};
</script>
<style lang="sass" scoped>
.regionSelect
  &.el-cascader
    width: 100%
</style>

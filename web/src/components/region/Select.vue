<template>
  <el-cascader
    class="region"
    v-if="inited"
    clearable
    :props="props"
    @change="handleChange"
    v-model="region"
    :show-all-levels="$props.showAllLevels"
  ></el-cascader>
</template>
<script>
import { mapActions, mapState } from "vuex";

const countryLevel = 0;
const chinaCode = "CN";

export default {
  name: "Region",
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
    return {
      inited: false,
      region: this.$props.value || "",
      query: {
        limit: 100,
        status: 1,
        fields: "code,name"
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
      try {
        const { regions } = await this.listRegion({
          params: Object.assign(
            {
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
.region
  .el-cascader
    width: 100%
</style>

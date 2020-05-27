<script>
import { CONFIG_EDITE_MODE } from "@/constants/route";
export default {
  name: "BaseTable",
  computed: {
    currentPage() {
      const { offset, limit } = this.query;
      return Math.floor(offset / limit) + 1;
    },
    editMode() {
      return this.$route.query.mode === CONFIG_EDITE_MODE;
    }
  },
  methods: {
    handleCurrentChange(page) {
      this.query.offset = (page - 1) * this.query.limit;
      this.fetch();
    },
    handleSizeChange(pageSize) {
      this.query.limit = pageSize;
      this.query.offset = 0;
      this.fetch();
    },
    add() {
      this.$router.push({
        query: {
          mode: CONFIG_EDITE_MODE
        }
      });
    },
    modify(item) {
      this.$router.push({
        query: {
          mode: CONFIG_EDITE_MODE,
          id: item.id
        }
      });
    }
  },
  watch: {
    $route() {
      if (!this.editMode) {
        this.fetch();
      }
    }
  },
  beforeMount() {
    if (!this.editMode) {
      this.fetch();
    }
  }
};
</script>

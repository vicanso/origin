import request from "@/helpers/request";
import { REGIONS, REGIONS_LIST_CATEGORIES, REGIONS_ID } from "@/constants/url";
import { attachStatusDesc, attachUpdatedAtDesc } from "@/store/modules/common";
import { findByID } from "@/helpers/util";

const state = {
  categoriesListProcessing: false,
  categories: null,
  processing: false,
  updateProcessing: false,
  list: {
    data: null,
    count: -1,
    country: {
      data: null,
      count: -1
    },
    province: {
      data: null,
      count: -1
    },
    city: {
      data: null,
      count: -1
    },
    area: {
      data: null,
      count: -1
    },
    street: {
      data: null,
      count: -1
    }
  }
};

const prefix = "region";
const mutationRegionList = `${prefix}.list`;
const mutationRegionListProcessing = `${mutationRegionList}.processing`;

const mutationRegionListCategory = `${prefix}.list.category`;
const mutationRegionListCategoryProcessing = `${mutationRegionListCategory}.processing`;

const mutationRegionUpdate = `${prefix}.update`;
const mutationRegionUpdateProcessing = `${mutationRegionUpdate}.processing`;

// listRegionCategory 获取地区分类
async function listRegionCategory({ commit }) {
  if (state.categoies) {
    return {
      categoies: state.categories
    };
  }
  commit(mutationRegionListCategoryProcessing, true);
  try {
    const { data } = await request.get(REGIONS_LIST_CATEGORIES);
    commit(mutationRegionListCategory, data);
    return data;
  } finally {
    commit(mutationRegionListCategoryProcessing, false);
  }
}

export default {
  state,
  mutations: {
    [mutationRegionListProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationRegionList](state, { categoy, regions, count }) {
      let data = state.list;
      if (categoy) {
        data = state[categoy];
      }
      if (count >= 0) {
        data.count = count;
      }
      regions.forEach(item => {
        attachUpdatedAtDesc(item);
        attachStatusDesc(item);
      });
      data.data = regions;
    },
    [mutationRegionListCategoryProcessing](state, processing) {
      state.categoriesListProcessing = processing;
    },
    [mutationRegionListCategory](state, { categories }) {
      state.categories = categories;
    },
    [mutationRegionUpdateProcessing](state, processing) {
      state.updateProcessing = processing;
    },
    [mutationRegionUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
      }
    }
  },
  actions: {
    async listRegion({ commit }, { params, categoy }) {
      commit(mutationRegionListProcessing, true);
      try {
        const { data } = await request.get(REGIONS, {
          params
        });
        data.categoy = categoy;
        commit(mutationRegionList, data);
        return data;
      } finally {
        commit(mutationRegionListProcessing, false);
      }
    },
    listRegionCategory,
    async updateRegionByID({ commit }, { id, data }) {
      commit(mutationRegionUpdateProcessing, true);
      try {
        const url = REGIONS_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationRegionUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationRegionUpdateProcessing, false);
      }
    },
    async getRegionByID({ commit }, id) {
      const found = findByID(state.list.data, id);
      if (found) {
        return found;
      }
      commit(mutationRegionListProcessing, true);
      try {
        const url = REGIONS_ID.replace(":id", id);
        const { data } = await request.get(url);
        return data;
      } finally {
        commit(mutationRegionListProcessing, false);
      }
    }
  }
};

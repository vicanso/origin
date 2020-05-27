import request from "@/helpers/request";

import { BRANDS_LIST_STATUS, BRANDS, BRANDS_ID } from "@/constants/url";
import { formatDate, addNoCacheQueryParam } from "@/helpers/util";

const mutationBrandListStatus = "brand.list.status";
const mutationBrandListStatusProcessing = `${mutationBrandListStatus}.processing`;
const mutationBrandProcessing = "brand.processing";
const mutationBrandList = "brand.list";
const mutationBrandUpdate = "brand.update";

const state = {
  statusesListProcessing: false,
  statuses: null,
  processing: false,
  list: {
    data: null,
    count: -1
  }
};

// listBrandStatus 获取品牌状态列表
async function listBrandStatus({ commit }) {
  if (state.statuses) {
    return;
  }
  commit(mutationBrandListStatusProcessing, true);
  try {
    const { data } = await request.get(BRANDS_LIST_STATUS);
    commit(mutationBrandListStatus, data);
  } finally {
    commit(mutationBrandListStatusProcessing, false);
  }
}

function updateStatusDesc(item) {
  state.statuses.forEach(status => {
    if (item.status === status.value) {
      item.statusDesc = status.name;
    }
  });
}

export default {
  state,
  mutations: {
    [mutationBrandListStatusProcessing](state, processing) {
      state.statusesListProcessing = processing;
    },
    [mutationBrandListStatus](state, { statuses }) {
      state.statuses = statuses;
    },
    [mutationBrandProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationBrandList](state, { brands = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      brands.forEach(item => {
        item.updatedAtDesc = formatDate(item.updatedAt);
        updateStatusDesc(item);
      });
      state.list.data = brands;
    },
    [mutationBrandUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      state.list.data.forEach(item => {
        if (item.id === id) {
          Object.assign(item, data);
          updateStatusDesc(item);
        }
      });
    }
  },
  actions: {
    // addBrand 添加品牌信息
    async addBrand({ commit }, brand) {
      commit(mutationBrandProcessing, true);
      try {
        const { data } = await request.post(BRANDS, brand);
        return data;
      } finally {
        commit(mutationBrandProcessing, false);
      }
    },
    listBrandStatus,
    // listBrand 获取品牌
    async listBrand({ commit }, params) {
      commit(mutationBrandProcessing, true);
      try {
        await listBrandStatus({ commit });
        const { data } = await request.get(BRANDS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationBrandList, data);
      } finally {
        commit(mutationBrandProcessing, false);
      }
    },
    // getBrandByID 通过id获取brand信息
    async getBrandByID({ commit }, id) {
      if (state.list.data) {
        let found = null;
        state.list.data.forEach(item => {
          if (item.id === id) {
            found = item;
          }
        });
        if (found) {
          return found;
        }
      }
      commit(mutationBrandProcessing, true);
      try {
        const url = BRANDS_ID.replace(":id", id);
        const { data } = await request.get(url, {
          params: addNoCacheQueryParam()
        });
        return data;
      } finally {
        commit(mutationBrandProcessing, false);
      }
    },
    // updateBrandByID 通过id更新brand信息
    async updateBrandByID({ commit }, { id, data }) {
      if (!data || Object.keys(data).length === 0) {
        return;
      }
      commit(mutationBrandProcessing, true);
      try {
        const url = BRANDS_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationBrandUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationBrandProcessing, false);
      }
    }
  }
};

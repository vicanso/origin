import request from "@/helpers/request";

import { attachUpdatedAtDesc, attachStatusDesc } from "@/store/modules/common";
import {
  addNoCacheQueryParam,
  toUploadFiles,
  findByID,
  formatDate
} from "@/helpers/util";
import {
  ADVERTISEMENTS,
  ADVERTISEMENTS_ID,
  ADVERTISEMENT_CATEGORIES
} from "@/constants/url";

const prefix = "advertisement";
const mutationAdvertisementProcessing = `${prefix}.processing`;
const mutationAdvertisementList = `${prefix}.list`;
const mutationAdvertisementUpdate = `${prefix}.update`;
const mutationAdvertisementCategories = `${prefix}.categories`;
const mutationAdvertisementCategoriesProcessing = `${mutationAdvertisementCategories}.processing`;

const state = {
  processing: false,
  list: {
    data: null,
    count: -1
  },
  processingCategories: false,
  categories: null
};

function enhanceAdvertisementInfo(item) {
  item.startedAtDesc = formatDate(item.startedAt);
  item.endedAtDesc = formatDate(item.endedAt);
  attachUpdatedAtDesc(item);
  attachStatusDesc(item);
  item.files = toUploadFiles(item.pic);
}

export default {
  state,
  mutations: {
    [mutationAdvertisementProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationAdvertisementList](state, { advertisements = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      advertisements.forEach(enhanceAdvertisementInfo);
      state.list.data = advertisements;
    },
    [mutationAdvertisementUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
        enhanceAdvertisementInfo(found);
      }
    },
    [mutationAdvertisementCategoriesProcessing](state, processing) {
      state.processingCategories = processing;
    },
    [mutationAdvertisementCategories](state, { categories }) {
      state.categories = categories;
    }
  },
  actions: {
    // addAdvertisement 添加广告信息
    async addAdvertisement({ commit }, advertisement) {
      if (advertisement.files) {
        advertisement.pic = advertisement.files[0].url;
        delete advertisement.files;
      }
      commit(mutationAdvertisementProcessing, true);
      try {
        const { data } = await request.post(ADVERTISEMENTS, advertisement);
        return data;
      } finally {
        commit(mutationAdvertisementProcessing, false);
      }
    },
    // listAdvertisement 获取广告
    async listAdvertisement({ commit }, params) {
      commit(mutationAdvertisementProcessing, true);
      try {
        const { data } = await request.get(ADVERTISEMENTS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationAdvertisementList, data);
        return data;
      } finally {
        commit(mutationAdvertisementProcessing, false);
      }
    },
    // getAdvertisementByID 通过id获取advertisement信息
    async getAdvertisementByID({ commit }, id) {
      const found = findByID(state.list.data, id);
      if (found) {
        return found;
      }
      commit(mutationAdvertisementProcessing, true);
      try {
        const url = ADVERTISEMENTS_ID.replace(":id", id);
        const { data } = await request.get(url, {
          params: addNoCacheQueryParam()
        });
        enhanceAdvertisementInfo(data);
        return data;
      } finally {
        commit(mutationAdvertisementProcessing, false);
      }
    },
    // updateAdvertisementByID 通过id更新advertisement信息
    async updateAdvertisementByID({ commit }, { id, data }) {
      if (data.files) {
        data.pic = data.files[0].url;
        delete data.files;
      }
      if (!data || Object.keys(data).length === 0) {
        return;
      }
      commit(mutationAdvertisementProcessing, true);
      try {
        const url = ADVERTISEMENTS_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationAdvertisementUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationAdvertisementProcessing, false);
      }
    },
    async listAdvertisementsCategory({ commit }) {
      if (state.categories) {
        return {
          categories: state.categories
        };
      }
      commit(mutationAdvertisementCategoriesProcessing, true);
      try {
        const { data } = await request.get(ADVERTISEMENT_CATEGORIES, {
          params: addNoCacheQueryParam()
        });
        commit(mutationAdvertisementCategories, data);
        return data;
      } finally {
        commit(mutationAdvertisementCategoriesProcessing, false);
      }
    }
  }
};

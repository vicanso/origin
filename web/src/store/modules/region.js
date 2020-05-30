import request from "@/helpers/request";
import {
  REGIONS,
  COMMONS_STATUSES,
  REGIONS_LIST_CATEGORIES,
  REGIONS_ID
} from "@/constants/url";
import { formatDate } from "@/helpers/util";

const state = {
  statusesListProcessing: false,
  statuses: null,
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

const mutationRegionListStatus = `${prefix}.list.status`;
const mutationRegionListStatusProcessing = `${mutationRegionListStatus}.processing`;

const mutationRegionListCategory = `${prefix}.list.category`;
const mutationRegionListCategoryProcessing = `${mutationRegionListCategory}.processing`;

const mutationRegionUpdate = `${prefix}.update`;
const mutationRegionUpdateProcessing = `${mutationRegionUpdate}.processing`;

// listRegionStatus 获取地区状态
async function listRegionStatus({ commit }) {
  if (state.statuses) {
    return {
      statuses: state.statuses
    };
  }
  commit(mutationRegionListStatusProcessing, true);
  try {
    const { data } = await request.get(COMMONS_STATUSES);
    commit(mutationRegionListStatus, data);
    return data;
  } finally {
    commit(mutationRegionListStatusProcessing, false);
  }
}

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
        item.updatedAtDesc = formatDate(item.updatedAt);
        updateStatusDesc(item);
      });
      data.data = regions;
    },
    [mutationRegionListStatusProcessing](state, processing) {
      state.statusesListProcessing = processing;
    },
    [mutationRegionListStatus](state, { statuses }) {
      state.statuses = statuses;
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
      state.list.data.forEach(item => {
        if (item.id === id) {
          Object.assign(item, data);
        }
      });
    }
  },
  actions: {
    async listRegion({ commit }, { params, categoy }) {
      commit(mutationRegionListProcessing, true);
      try {
        await listRegionStatus({ commit });
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
    listRegionStatus,
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

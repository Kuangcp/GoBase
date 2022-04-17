window.api = {
    account: {
        listAll: "/account/listAccount",
    },
    user: {
        listAll: "/user/listUser",
    },
    category: {
        tree: "/category/listCategoryTree"
    },
    record: {
        create: "/record/createRecord",
        balance: "/record/calBalance",
        list: "/record/listRecord",
        byCategory: "/record/category",
        categoryDetail: "/record/categoryDetail"
    },
    loan: {
        create: "/loan/create",
        query: "/loan/query"
    },
    report: {
        categoryPeriod: "/report/categoryPeriod",
        balanceReport: "/report/balanceReport"
    }
}
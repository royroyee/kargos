
export const Routes = {
    // pages
    DashboardOverview: { path: "/" },

    // nodes
    NodesOverview: { path: "/nodes/overview/" },
    NodesDetail: { path: "/nodes/detail/:page"},

    // events
    EventsOverview: { path: "/events/overview/" },

    // resources overview
    PodsOverview: { path: "/resources/pods/overview/" },
    PodDetail: { path: "/resources/pods/detail/:namespace/:page"},

    WorkloadOverview: { path: "/workload/overview/"},
    NotFound: { path: "/404" },
    ServerError: { path: "/500" },
};
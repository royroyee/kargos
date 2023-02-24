import React, { useState, useEffect } from 'react';
import { Route, Switch, Redirect } from "react-router-dom";
import { Routes } from "../routes";

// single pages
import DashboardOverview from "./dashboard/DashboardOverview";

// node information pages
import NodesOverview from "./nodes/NodesOverview";
import NodesDetail from "./nodes/NodesDetail";

// events overview pages
import EventsOverview from "./events/EventsOverview";

// workload overview page
import WorkloadOverview from './workload/WorkloadOverview';

// resource pages
import PodDetail from './resources/PodDetail';

// error pages
import NotFoundPage from "./etc/NotFound";
import ServerError from "./etc/ServerError";

// components
import Sidebar from "../components/Sidebar";

const RouteWithLoader = ({ component: Component, ...rest }) => {
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setLoaded(true), 1000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <Route {...rest} render={props => ( <> <Component {...props} /> </> ) } />
  );
};

const RouteWithSidebar = ({ component: Component, ...rest }) => {
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setLoaded(true), 1000);
    return () => clearTimeout(timer);
  }, []);

  const localStorageIsSettingsVisible = () => {
    return localStorage.getItem('settingsVisible') === 'false' ? false : true
  }

  const [showSettings, setShowSettings] = useState(localStorageIsSettingsVisible);

  const toggleSettings = () => {
    setShowSettings(!showSettings);
    localStorage.setItem('settingsVisible', !showSettings);
  }

  return (
    <Route {...rest} render={props => (
      <>
        <Sidebar />

        <main className="content">
          <Component {...props} />
        </main>
      </>
    )}
    />
  );
};

export default () => (
  <Switch>
    <RouteWithLoader exact path={Routes.NotFound.path} component={NotFoundPage} />
    <RouteWithLoader exact path={Routes.ServerError.path} component={ServerError} />

    {/* pages */}
    <RouteWithSidebar exact path={Routes.DashboardOverview.path} component={DashboardOverview} />
    <RouteWithSidebar exact path={Routes.NodesOverview.path} component={NodesOverview} />
    <RouteWithSidebar exact path={Routes.NodesDetail.path} component={NodesDetail} />

    <RouteWithSidebar exact path={Routes.EventsOverview.path} component={EventsOverview} />
    <RouteWithSidebar exact path={Routes.WorkloadOverview.path} component={WorkloadOverview} />

    <RouteWithSidebar exact path={Routes.PodDetail.path} component={PodDetail} />

    <Redirect to={Routes.NotFound.path} />
  </Switch>
);

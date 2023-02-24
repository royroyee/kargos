
import React from "react";
import { Col, Row } from '@themesberg/react-bootstrap';
import { DashboardPage } from "../../components/Dashboard";

export default () => {
  return (
    <>
      <article>
            <Row className="d-flex flex-wrap flex-md-nowrap py-4">
              <Col className="d-block mb-4 mb-md-0">
                <h1 className="h2">Overview</h1>
                <p className="mb-0">
                  Kubernetes Cluster Overview
                </p>
              </Col>
            </Row>
      </article>
    
      <DashboardPage></DashboardPage>
    </>
  );
};

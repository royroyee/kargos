import React from "react";
import { Col, Row} from '@themesberg/react-bootstrap';

import { NodesOverviewPage } from "../../components/Nodes";

export default () => {
  return (
    <article>
          <Row className="d-flex flex-wrap flex-md-nowrap py-4">
            <Col className="d-block mb-4 mb-md-0">
              <h1 className="h2">Nodes</h1>
              <p className="mb-0">
              Node status in Kubernetes Cluster.
              </p>
            </Col>
          </Row>
       <NodesOverviewPage></NodesOverviewPage>
    </article>
  );
};
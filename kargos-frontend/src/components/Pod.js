import { Col, Row, Card, Table } from '@themesberg/react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCompactDisc, faServer } from '@fortawesome/free-solid-svg-icons';

import './css/CodeBlock.css'
import podDetail from "../fakeData/DUMMY_POD_INFO"

import React, { useState, useEffect } from 'react';
import Chartist from "react-chartist";
import ChartistTooltip from 'chartist-plugin-tooltips-updated';

import AccordionComponent from '../components/AccordionComponent'


/**
 * Generate pod detail page.
 * @param {props} props The props
 * @returns A JSX component that represents pod detail page.
 */
export const PodDetailPage = (props) => {
    const { page, namespace } = props;
    console.log (page + namespace)
    return (
      <>
          <article>
            <Row className="d-flex flex-wrap flex-md-nowrap py-4">
              <Col className="d-block mb-4 mb-md-0">
              <h1 className="h2">Pod Information</h1>
              <p className="mb-0">
                  {page}
              </p>
              <small>
                Namespace : {namespace}
              </small>
              </Col>
            </Row>
          </article>
          <LastStatus name={page}></LastStatus>
          <Row className="justify-content-md-center">
            <Col xs={12} xl={12} className="mb-4">
              <PodLogWidget name={page} namespace={namespace}/>
            </Col>
          </Row>
          <Row className="justify-content-md-center">
            <Col xs={12} xl={6} className="mb-4">
              <PodInfoWidget name={page} namespace={namespace}/>
            </Col>
            <Col xs={12} xl={6} className="mb-4">
              <PodInfoWidget name={page} namespace={namespace}/>
            </Col>
          </Row>
          <Row className="justify-content-md-center">
            <PodContainersTable></PodContainersTable>
          </Row>
      </>
    );
}

/**
 * Generate last 24 hours graph for cpu, ram, disk and network.
 * @returns A JSX Component that represents last 24 hours graph.
 */
const LastStatus = (props) => {
  const { name } = props;
  const [status, setStatus] = useState({ "cpu_usage": [], "ram_usage": [], "network_usage": [], "disk_usage": [] });

  // Retrieve data from REST API.
  const getStatus = () => {
    const url = "/pod/usage/" + name;
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    fetch(url, requestOptions)
        .then(response => response.text())
        .then(result => {
          setStatus(JSON.parse(result));
        })
        .catch(error => console.log('error', error));
  }

  useEffect(() => {
      getStatus();
  }, []);

  // Generate BarWidget using data retrieved.
  return ( 
      <>
      <Row className="justify-content-md-center">
          <Col xs={12} xl={6} className="mb-4">
              <LastStatusGraph data={status.cpu_usage} title={"CPU Usage"} metric={"Mi"}></LastStatusGraph>
          </Col>
          <Col xs={12} xl={6} className="mb-4">
              <LastStatusGraph data={status.ram_usage} title={"RAM Usage"} metric={"MB"}></LastStatusGraph>
          </Col>
      </Row>
      <Row className="justify-content-md-center">
          <Col xs={12} xl={6} className="mb-4">
              <LastStatusGraph data={status.network_usage} title={"Network Usage"} metric={"%"}></LastStatusGraph>
          </Col>
          <Col xs={12} xl={6} className="mb-4">
              <LastStatusGraph data={status.disk_usage} title={"Disk Usage"} metric={"%"}></LastStatusGraph>
          </Col>
      </Row>
      </>
  );
}

/**
* Generate last 24 hours graph.
* @param {props} props The props
* @returns A JSX Component that shows last status graph for 24 hours.
*/
const LastStatusGraph = (props) => {
  const { data, title, metric } = props;
  return (
      <Card className="shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              {title}
            </h5>
            <h3>{data === null || data === undefined ? "UNKNOWN" : data[data.length - 1] + metric}</h3>
            <small>Last 24 Hours</small>
          </div>
        </Card.Header>
        <Card.Body className="p-2">
          <Valuechart data={data === null || data === undefined ? [] : data}></Valuechart>
        </Card.Body>
      </Card>
    );
}

/**
* Generate value chart from data.
* @param {props} props The props
* @returns A JSX Component that shows graph in series.
*/
export const Valuechart = (props) => {
  const { data } = props; 
  var graphValues = data;

  // If data was null or undefined, set it 0.
  if (data === null || data === undefined) {
      data = new Array(24).fill(0);
  } else if (data.length != 24) { // When data was not enough for last 24 hours, fill rest with 0.
      var original = data;
      var newData = new Array(24).fill(0);

      for (let i = 0; i < original.length; i++) {
          newData[i] = original[i];
      }
      graphValues = newData;
  }

  // Add labels
  var label = []
  for (var i = 0; i < data.length; i++) {
    label.push(24 - i)
  }

  const graphData = {
    labels: label,
    series: [graphValues]
  }

  const options = {
    low: 0,
    showArea: true,
    fullWidth: true,
    axisX: {
      position: 'end',
      showGrid: true
    },
    axisY: {
      // On the y-axis start means left and end means right
      showGrid: false,
      showLabel: false,
    }
  };

  const plugins = [
    ChartistTooltip()
  ]

  return (
    <Chartist data={graphData} options={{...options, plugins}} type="Line" className="ct-series-g ct-double-octave" />
  );
}

/**
 * A function that generates a log of pod in code block style.
 * @param {props} props The props to use when generating this node info table.
 * @returns A JSX Component that represents the node log code block.
 */
const PodLogWidget = (props) => {
  const { name, namespace } = props;
  const [log, setLog] = useState([]);

  // Retrieve data from REST API.
  const getLog = () => {
    const url = "/pod/logs/" + namespace + "/" + name;
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    console.log(url)

    fetch(url, requestOptions)
        .then(response => response.text())
        .then(result => {
          setLog(JSON.parse(result));
        })
        .catch(error => console.log('error', error));
  }

  useEffect(() => {
    getLog();
  }, []);

    var outText = "";
    for (let i = log.length - 1; i > log.length - 12; i--) {
        if (log[i] === undefined || log[i] === null) {
          outText = outText + "\n";
        } else {
          outText = outText + log[i] + '\n';
        }
    }

    return (
      <Card className="shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              Pod Log
            </h5>
          </div>
        </Card.Header>
        <Card.Body className="p-2">
        <div className="code-block">
            <pre>
                <code>{outText}</code>
            </pre>
        </div>
        </Card.Body>
      </Card>
    );
}

/**
 * Generate pod info widget.
 * @param {props} props The props
 * @returns A JSX Component that represents pod info widget
 */
const PodInfoWidget = (props) => {
    const { name, namespace } = props;
    return (
        <Card className="shadow-sm">
        <PodInfoTable name={name} namespace={namespace}></PodInfoTable>
        </Card>
    );
}

/**
 * A function that shows pod information
 * @param {*} page the page name
 * @returns Table consisting of all pod informations.
 */
export const PodInfoTable = (props) => {

    const { name, namespace } = props;
    const [info, setInfo] = useState({});

    // Retrieve data from REST API.
    const getInfo = () => {
      // @todo add namespace filtering in REST API when it is available.
      const url = "/pod/info/" + name;
      console.log(url)
      var requestOptions = {
          method: 'GET',
          redirect: 'follow'
      };

      console.log(url)

      fetch(url, requestOptions)
          .then(response => response.text())
          .then(result => {
            setInfo(JSON.parse(result));
          })
          .catch(error => console.log('error', error));
    }

    useEffect(() => {
      getInfo();
    }, []);


    const TableRow = (props) => {
      const { type, value } = props;
      // TODO ADD LINK LIKE CARD TO THE NODE's NAME and delete action
      return (
        <tr>
          <td>
            <span className="fw-normal">
              {type}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              { Array.isArray(value) ? value.join(", ") : value}
            </span>
          </td>
        </tr>
      );
    };
    
    // Generate Table rows by iterating key and values.
    var items = [];
    const keys = Object.keys(info);

    // Iterate over keys and generate table rows.
    for (const key of keys) {
        items.push(<TableRow type={key} value={info[key]} />);
    }
  
    console.log(items)

    return (
      <Card border="light" className="table-wrapper table-responsive shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              Pod Information
            </h5>
          </div>
        </Card.Header>
  
        <Card.Body className="pt-0">
          <Table hover className="user-table align-items-center">
            <thead>
              <tr>
                <th className="border-bottom">Type</th>
                <th className="border-bottom">Value</th>
              </tr>
            </thead>
            <tbody>
              {items}
            </tbody>
          </Table>
        </Card.Body>
      </Card>
    );
  };

/**
 * Generate pod container dropdown.
 * @todo add REST API.
 * @param {props} props The props.
 * @returns A JSX Component that represents pod container dropdown.
 */
export const PodContainerListWidget = (props) => {
return (
    <Card className="shadow-sm">
    <PodContainersTable></PodContainersTable>
    </Card>
);
}

/**
 * A function that shows pod information
 * @param {*} page the page name
 * @returns Table consisting of all pod informations.
 */
export const PodContainersTable = () => {
    const ProcessRow = (process) => {
      return (
        <tr>
          <td>
            <span className="fw-normal">
              {process.name}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {process.status}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {process.PID}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {process.cpu}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {process.ram}
            </span>
          </td>
        </tr>
      );
    };
  
    // Generate Cards for each of those containers.
    const ContainersCards = (container) => {
      return (
        <Row className="justify-content-md-center">
          <Card border="light" className="table-wrapper table-responsive shadow-sm">
            <Card.Header className="d-flex flex-row align-items-center flex-0">
              <div className="d-block">
                <h5 className="fw-normal mb-2">
                  ContainerID : {container.id}
                </h5>
                <small>
                  <FontAwesomeIcon icon={faCompactDisc}></FontAwesomeIcon> <b>Image</b> : {container.image}
                </small>
                <br></br>
                <small>
                  <FontAwesomeIcon icon={faServer}></FontAwesomeIcon> <b>Node</b> : {container.node}
                </small>
              </div>
            </Card.Header>
  
            <Card.Body className="pt-0">
              <Table hover className="user-table align-items-center">
                <thead>
                  <tr>
                    <th style={{ width: '30%' }} className="border-bottom">Name</th>
                    <th style={{ width: '10%' }} className="border-bottom">Status</th>
                    <th style={{ width: '10%' }} className="border-bottom">PID</th>
                    <th style={{ width: '25%' }} className="border-bottom">CPU Usage</th>
                    <th style={{ width: '25%' }} className="border-bottom">RAM Usage</th>
                  </tr>
                </thead>
                <tbody>
                  {container.processes.map(t => <ProcessRow key={`processes-${t.index}`} {...t} />)}
                </tbody>
              </Table>
            </Card.Body>
          </Card>
        </Row>
      );
    }
  
    var accordionElements = [];
    for (var i = 0; i < podDetail.containers.length; i++) {
      var container = podDetail.containers[i];
      var name = container.id;
      accordionElements.push({
        id: name,
        eventKey: name,
        title: name,
        description: 
        <Card border="light" className="table-wrapper table-responsive shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              ContainerID : {container.id}
            </h5>
            <small>
              <FontAwesomeIcon icon={faCompactDisc}></FontAwesomeIcon> <b>Image</b> : {container.image}
            </small>
            <br></br>
            <small>
              <FontAwesomeIcon icon={faServer}></FontAwesomeIcon> <b>Node</b> : {container.node}
            </small>
          </div>
        </Card.Header>
  
        <Card.Body className="pt-0">
          <Table hover className="user-table align-items-center">
            <thead>
              <tr>
                <th style={{ width: '30%' }} className="border-bottom">Name</th>
                <th style={{ width: '10%' }} className="border-bottom">Status</th>
                <th style={{ width: '10%' }} className="border-bottom">PID</th>
                <th style={{ width: '25%' }} className="border-bottom">CPU Usage</th>
                <th style={{ width: '25%' }} className="border-bottom">RAM Usage</th>
              </tr>
            </thead>
            <tbody>
              {container.processes.map(t => <ProcessRow key={`processes-${t.index}`} {...t} />)}
            </tbody>
          </Table>
        </Card.Body>
      </Card>
        });
    }
  
    return (
      <Card border="light" className="table-wrapper table-responsive shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              Running Containers
            </h5>
          </div>
        </Card.Header>
        <Card.Body className="pt-0">
          <AccordionComponent
            defaultKey = "none"
            data = {
              accordionElements
            }
          />
        </Card.Body>
      </Card>
    );
  };
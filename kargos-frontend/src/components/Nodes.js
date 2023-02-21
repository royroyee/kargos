
import { useState, useEffect } from "react";
import { Col, Row, Card, Table, Pagination } from '@themesberg/react-bootstrap';
import { Link } from 'react-router-dom';
import { ResourceChart } from "./Charts";

import Chartist from "react-chartist";
import ChartistTooltip from 'chartist-plugin-tooltips-updated';

import nodeInfo from "../fakeData/DUMMY_NODE_INFO"
import './css/CodeBlock.css'

/**
 * Generate overall page for the Nodes.
 * @returns A JSX Component that shows nodes table, tab and pagenation.
 */
export const NodesOverviewPage = () => {
    return <NodesTableSection></NodesTableSection>
}

/**
 * Generate detailed page for a specific node.
 * @returns A JSX Component that shows node page.
 */
export const NodesDetailPage = (props) => {
    const { page } = props;
    const name = page.page;

    // Get node information using REST API from backend.
    const [data, setData] = useState([]);
    const getData = () => {
        const url = "/node/info/" + name;

        var requestOptions = {
            method: 'GET',
            redirect: 'follow'
        };
    
        fetch(url, requestOptions)
            .then(response => response.text())
            .then(result => {
              setData(JSON.parse(result));
            })
            .catch(error => console.log('error', error));
    }

    useEffect(() => {
        getData();
    }, []);

    return (
        <article>
            <Row className="d-flex flex-wrap flex-md-nowrap py-4">
                <Col className="d-block mb-4 mb-md-0">
                    <h1 className="h2">{name}</h1>
                </Col>
            </Row>
            <LastStatus name={name}></LastStatus>
            <Row className="justify-content-md-center">
                <Col xs={12} xl={12} className="mb-4">
                    <NodeInfoTable data={data}/>
                </Col>
            </Row>
            <Row className="justify-content-md-center">
                <Col xs={12} xl={12} className="mb-4">
                    <NodeLogWidget page={page.page}/>
                </Col>
            </Row>
        </article>
    );
}

/**
 * A function that generates a table in type: value order for node information.
 * @param {props} props The props to use when generating this node info table.
 * @returns A JSX Component that represents the node information table.
 */
const NodeInfoTable = (props) => {
    const { data } = props;

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
              {value}
            </span>
          </td>
        </tr>
      );
    };
    
    // Generate Table rows by iterating key and values.
    var items = [];
    const keys = Object.keys(data);

    // Iterate over keys and generate table rows.
    for (const key of keys) {
      if (data[key] != true) 
        items.push(<TableRow type={key} value={data[key]} />);
      else 
        items.push(<TableRow type={key} value="Ready" />);
    }

    return (
    <>
        <Card border="light" className="table-wrapper table-responsive shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
            <div className="d-block">
                <h5 className="fw-normal mb-2">
                Node Information
                </h5>
            </div>
            </Card.Header>
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
        </Card>
    </>
    );
}

/**
 * A function that generates a log of kubelet in code block style.
 * @param {props} props The props to use when generating this node info table.
 * @returns A JSX Component that represents the node log code block.
 */
const NodeLogWidget = () => {
    return (
      <Card className="shadow-sm">
        <Card.Header className="d-flex flex-row align-items-center flex-0">
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              Kubelet Log
            </h5>
          </div>
        </Card.Header>
        <Card.Body className="p-2">
        <div className="code-block">
            <pre>
                <code>{nodeInfo.nodeLog}</code>
            </pre>
        </div>
        </Card.Body>
      </Card>
    );
  }

/**
 * Generates the list of nodes in a specific page.
 * @param {page} page The page information.
 * @returns The Card and a Table element inside that table representing nodes.
 */
export const NodesTableSection = () => {
    const [data, setData] = useState([]);
    const [page, setPage] = useState(1);

    /**
     * This function will handle new page click from pagination. Also will trigger the code to re-render the table.
     * @param {int} newPage The new page to set.
     */
    function handlePageinationClick(newPage) {
        setPage(newPage);
        updateTable(newPage); // Need to send selection since the setPage might not be updated yet.
    }

    /**
     * Generate Pageination for Nodes.
     * @todo add support for relative pagenation. For example, when we have 100 tabs, we need to have 25 ~ 35 printed out in the screen. But for now, it does not.
     * @returns JSX component that implements pageination.
     */
    const NodesPagination = () => {
        const [count, setCount] = useState([]);

        const getCount = () => {
            const url = "/nodes/count";
            var requestOptions = {
                method: 'GET',
                redirect: 'follow'
            };
        
            fetch(url, requestOptions)
                .then(response => response.text())
                .then(result => {
                    setCount(JSON.parse(result));
                })
                .catch(error => console.log('error', error));
        }
            
        useEffect(() => {
            getCount();
        }, []);
        
        var nodesCount = count['count']
        
        const items = [];
        var totalPages = Math.ceil(nodesCount / 10) + 1;      
        var prevDisabled = (page == 1);
        var nextDisabled = (page == totalPages - 1) || totalPages == 1;

        items.push(
            <Pagination.Prev disabled={prevDisabled} onClick={() => handlePageinationClick((page - 1))}>
            Previous
            </Pagination.Prev>
        );
        
        for (let i = 1; i < totalPages; i++) {
            if (i != Number(page)) {
            items.push(<Pagination.Item onClick={() => handlePageinationClick(i)}>
                {i}
            </Pagination.Item>);
            } else {
            items.push(<Pagination.Item active>
                {i}
            </Pagination.Item>);
            }
        }
        
        items.push(
            <Pagination.Next disabled={nextDisabled} onClick={() => handlePageinationClick((page + 1))}>
            Next
            </Pagination.Next>
        );
        
        return (
            <Pagination className="mb-2 mb-lg-0">
                {items}
            </Pagination>
        );
    }

    /**
     * This function will update table from data using setData and refresh the table for the user.
     * This uses REST API /nodes/?page= API from backend.
     * @param {Number} argPage The page
     */
    function updateTable (argPage) {
        const url = "/nodes/" + "?page=" + argPage + "&per_page=10";
        console.log(url)
        var requestOptions = {
            method: 'GET',
            redirect: 'follow'
        };
    
        fetch(url, requestOptions)
            .then(response => response.text())
            .then(result => {
                setData(JSON.parse(result));
            })
            .catch(error => console.log('error', error));
        
    }

    useEffect(() => {
        updateTable(page);
    }, []);

    /**
     * Generate table rows for each elements
     * @param {*} props The props that includes each data.
     * @returns A JSX Component that represents a set of Rows in table.
     */
    const TableRow = (props) => {
      const { name, cpu_usage, ram_usage, disk_allocated, network_usage, ip, status } = props;
      const statusVariant = status === "Ready" ? "success"
        : status === "Not Ready" ? "warning"
          : status === "No Connection" ? "danger" : "primary";
      return (
        <tr>
          <td>
            <Link to={"/nodes/detail/" + name}> {name} </Link>
          </td>
          <td>
            <span className="fw-normal">
              {cpu_usage}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {ram_usage}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {disk_allocated}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {network_usage}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {ip}
            </span>
          </td>
          <td>
            <span className={`fw-normal text-${statusVariant}`}>
              {status}
            </span>
          </td>
        </tr>
      );
    };

    /**
     * Generate a table for nodes.
     * @returns A JSX Component that represents nodes table.
     */
    const NodesTable = () => {
        return (
            <Table hover className="user-table align-items-center">
                <thead>
                <tr>
                    <th className="border-bottom">Node Name</th>
                    <th className="border-bottom">CPU Usage (%)</th>
                    <th className="border-bottom">RAM Usage (%)</th>
                    <th className="border-bottom">Disk Usage (%)</th>
                    <th className="border-bottom">Network Usage (%)</th>
                    <th className="border-bottom">IP</th>
                    <th className="border-bottom">Status</th>
                </tr>
                </thead>
                <tbody>
                { // If the retrieved data was null, set it empty table row.
                    data ? (
                    <>
                        {data.map(t => (<TableRow {...t} />
                        ))} 
                    </>
                    ) : (
                        <tr>
                        </tr>
                    )
                }
                </tbody>
            </Table>
        );
    }

    return (
    <Card border="light" className="table-wrapper table-responsive shadow-sm">
        <Card.Body className="pt-0">
            <NodesTable></NodesTable>
        </Card.Body>
        <Card.Footer className="px-3 border-0 d-lg-flex align-items-center justify-content-between">
            <NodesPagination></NodesPagination>
        </Card.Footer>
    </Card>
    );
  };


/**
 * Generate last 24 hours graph for cpu, ram, disk and network.
 * @returns A JSX Component that represents last 24 hours graph.
 */
const LastStatus = (props) => {
  const { name } = props;
  const [status, setStatus] = useState({ "cpu_usage": [], "ram_usage": [], "network_usage": [], "disk_usage": [] });

  // Retrieve data from REST API.
  const getStatus = () => {
    const url = "/node/usage/" + name;
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
              <LastStatusGraph data={status.cpu_usage} title={"CPU Usage"} metric={"%"}></LastStatusGraph>
          </Col>
          <Col xs={12} xl={6} className="mb-4">
              <LastStatusGraph data={status.ram_usage} title={"RAM Usage"} metric={"%"}></LastStatusGraph>
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
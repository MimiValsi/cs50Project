(async () => {
  const getJson = async () => {
    const response = await fetch("http://localhost:3005/jsonGraph");
    const data = await response.json();
    const js = await JSON.stringify(data);
    const jj = await JSON.parse(js)

    return jj;
  };

  const jsData = await getJson();

  let infoNb = [];
  for (let i = 0; i < jsData.length; i++) {
    infoNb.push(jsData[i].curatifs);
  }

  let sourceName = [];
  for (let i = 0; i < jsData.length; i++) {
    sourceName.push(jsData[i].name);
  }

  // for ESM environment, need to import modules as:
  // import bb, {bar} from "billboard.js";

  var chart = bb.generate({
    bindto: "#myPlot",

    data: {
      names: {
        data1: "Info",
        data2: "Sources"
      },
      columns: [
        ["Number of info", ...infoNb],
      ],
      type: "bar", // for ESM specify as: bar()
    },

    axis: {
      x: {
        type: "category",
        categories: [...sourceName],
        height: 50,
        tick: {
          rotate: 75,
          multiline: false,
        }
      }
    },

    size: {
      width: 1000,
      height: 400
    },

    padding: true,

    resize: true,

    zoom: {
      enabled: true,
      type: "drag"
    },

    legend: {
      position: "inset"
    },

    bar: {
      width: {
        ratio: 0.5
      }
    }
  });
})();

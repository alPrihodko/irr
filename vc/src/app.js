class StaticGlobals {
  //const HOST = "sasha123.ddns.ukrtel.net:1234";

  constructor() {
    console.log("constants initiates");
  };

  static GetHostName() {
    return "sasha123.ddns.ukrtel.net:1235";
  };

  static OFF () {
    return "Off";
  }

  static AUTO () {
    return "Auto";
  }

  static ON () {
    return "On";
  }

};

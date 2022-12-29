// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

class TableWidget extends View {

  constructor() {
    super('table');
    this.addClass('table-widget');
  }

}

class HeaderTableRowWidget extends View {

  constructor() {
    super('tr');
    this.addClass('table-header-widget');
  }

}

class AutoHeaderTableRowWidget extends HeaderTableRowWidget {

  constructor(stringArray) {
    super();
    for (const text of stringArray) {
      this.appendView(new TableHeaderWidget(text));
    }
  }

}

class TableHeaderWidget extends View {

  constructor(text, hoverText) {
    super('th');
    this.addClass('table-header-text-widget');
    this.appendView(new HoverTextView(text, hoverText));
  }

}

class AutoTableRowWidget extends View {

  constructor(stringArray) {
    super('tr');
    this.addClass('table-row-widget');
    for (const text of stringArray) {
      this.appendView(new TableDataTextWidget(text));
    }
  }

}

class HoverTextView extends View {

  constructor(text, hoverText) {
    super('span');
    if (hoverText !== undefined) {
      this.setAttribute('title', hoverText);
    }
    this.appendElement(document.createTextNode(text));
  }

}

class TableDataTextWidget extends View {

  constructor(text) {
    super('td');
    this.addClass('table-data-text-widget');
    this.appendElement(document.createTextNode(text));
  }

}

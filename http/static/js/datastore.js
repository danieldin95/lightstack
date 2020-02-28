import {DataStoreApi} from "./api/datastore.js";
import {CheckBoxTop} from "./com/utils.js";


export class DataStore {
    // nil
    constructor() {
        this.datastoreOn = new DataStoreOn();
        this.datastores = this.datastoreOn.uuids;

        // register buttons's  click.
        $("datastore-delete").on("click", this.datastores, function (e) {
            new DataStoreApi({uuids: e.data.store}).delete();
        });
    }

    create(data) {
        new DataStoreApi().create(data);
    }
}


export class DataStoreOn {

    constructor() {
        this.uuids = {store: []};

        let change = this.change;
        let record = this.uuids;

        new CheckBoxTop({
            one: "datastore-on-one input",
            all: "datastore-on-all input",
            change: function (e) {
                change(record, e);
            }
        });

        // disabled firstly.
        change(record, this.uuids);
    }

    change(record, from) {
        record.store = from.store;
        console.log(record.store);

        if (from.store.length == 0) {
            $("datastore-edit button").addClass('disabled');
            $("datastore-delete button").addClass('disabled');
        } else {
            $("datastore-edit button").removeClass('disabled');
            $("datastore-delete button").removeClass('disabled');
        }

        if (from.store.length != 1) {
            $("datastore-edit button").addClass('disabled');
        }
        else {
            $("datastore-edit button").removeClass('disabled');
        }
    }
}
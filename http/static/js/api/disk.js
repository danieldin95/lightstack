import {Api} from "./api.js"
import {AlertDanger, AlertSuccess, AlertWarn} from "../com/alert.js";


export class DiskApi extends Api {
    // {
    //   instance: 'uuid',
    //   uuids: [uuid],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);

        this.instance = props.instance;
    }

    url(instance, uuid) {
        if (uuid) {
            return `/api/instance/${instance}/disk/${uuid}`
        }
        return `/api/instance/${instance}/disk`
    }

    list(data, func) {
        let your = this;

        $.get(your.url(this.instance), {format: 'schema'}, function (resp, status) {
            func({data, resp});
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    create(data) {
        let your = this;

        $.post(your.url(this.instance), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`start disk ${data.message}`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }

    delete() {
        let your = this;

        this.uuids.forEach(function (uuid, index, err) {
            $.delete(your.url(your.instance, uuid), function (data, status) {
                $(your.tasks).append(AlertSuccess(`remove disk '${uuid}' ${data.message}`));
            }).fail(function (e) {
                $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
            });
        });
    }

    edit(data) {
        let your = this;

        $.post(your.url(your.instance, your.uuids[0]), JSON.stringify(data), function (data, status) {
            $(your.tasks).append(AlertSuccess(`edit disk '${data.name}' ${data.message}`));
        }).fail(function (e) {
            $(your.tasks).append(AlertDanger((`${this.type} ${this.url}: ${e.responseText}`)));
        });
    }
}

export class Widget {
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks;
        console.log("Widget", props);
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }
}
package structure

type JobExec struct {
	TaskScriptId           int    `json:"task_script_id"`
	TaskId                 int    `json:"task_id"`
	ScriptId               string `json:"script_id"`
	TaskGroupActionId      string `json:"TaskGroupActionId"`
	TaskActionId           string `json:"task_action_id"`
	TaskScriptActionId     string `json:"TaskScriptActionId"`
	TaskScriptHostActionId string `json:"TaskScriptHostActionId"`
	AgentIp                string `json:"agent_ip"`
	AgentPort              string `json:"agent_port"`
	ServerAddr             string `json:"server_addr"`
	RpcServerName          string `json:"rpc_server_name"`
	ProxyAddr              string `json:"proxy_addr"`
	ScriptName             string `json:"script_name"`
	Md5                    string `json:"md5"`
	Timeout                int    `json:"timeout"`
	Person                 string `json:"person"`
	RunUser                string `json:"run_user"`
	Command                string `json:"command"`
	ScriptContent          string `json:"script_content"`
	FileType               int8   `json:"file_type"`
	UploadPath             string `json:"upload_path"`
	Topic                  string `json:"topic"`
}

type CallBack struct {
	TaskScriptId           int    `json:"task_script_id"`
	TaskId                 int    `json:"task_id"`
	ScriptId               string `json:"script_id"`
	TaskGroupActionId      string `json:"TaskGroupActionId"`
	TaskActionId           string `json:"task_action_id"`
	TaskScriptActionId     string `json:"TaskScriptActionId"`
	TaskScriptHostActionId string `json:"TaskScriptHostActionId"`
	AgentIp                string `json:"agent_ip"`
	ServerAddr             string `json:"server_addr"`
	RpcServerName          string `json:"rpc_server_name"`
	ProxyAddr              string `json:"proxy_addr"`
	ScriptName             string `json:"script_name"`
	ScriptContent          string `json:"script_content"`
	Md5                    string `json:"md5"`
	Person                 string `json:"person"`
	RunUser                string `json:"run_user"`
	Command                string `json:"command"`
	StartTime              string `json:"start_time"`
	StopTime               string `json:"stop_time"`
	RunTime                int64  `json:"run_time"`
	Stdout                 string `json:"stdout"`
	Stderr                 string `json:"stderr"`
	Message                string `json:"message"`
	Status                 int8   `json:"status"`
	FileType               int8   `json:"file_type"`
	UploadPath             string `json:"upload_path"`
	Topic                  string `json:"topic"`
}

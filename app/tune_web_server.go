package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"fh6worker/internal/storage"
	"fh6worker/internal/telemetry"
)

type TuneWebServerStatus struct {
	Running    bool   `json:"running"`
	Port       int    `json:"port"`
	URL        string `json:"url"`
	LANAddress string `json:"lanAddress"`
	LastError  string `json:"lastError"`
}

type TuneWebServer struct {
	mu     sync.Mutex
	server *http.Server
	status TuneWebServerStatus
}

func NewTuneWebServer() *TuneWebServer {
	return &TuneWebServer{}
}

func (s *TuneWebServer) Start(port int) error {
	if port <= 0 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	s.mu.Lock()
	if s.server != nil {
		s.mu.Unlock()
		return errors.New("tune web server is already running")
	}
	address := preferredTuneWebAddress()
	listener, err := net.Listen("tcp", net.JoinHostPort(address, strconv.Itoa(port)))
	if err != nil {
		s.status = TuneWebServerStatus{Running: false, Port: port, LANAddress: address, LastError: err.Error()}
		s.mu.Unlock()
		return err
	}
	actualPort := listener.Addr().(*net.TCPAddr).Port
	mux := http.NewServeMux()
	mux.HandleFunc("/", tuneWebRootHandler)
	mux.HandleFunc("/tune", tuneWebPageHandler)
	mux.HandleFunc("/api/tune/generate", tuneWebGenerateHandler)
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	s.server = server
	s.status = TuneWebServerStatus{
		Running:    true,
		Port:       actualPort,
		LANAddress: address,
		URL:        fmt.Sprintf("http://%s/tune", net.JoinHostPort(address, strconv.Itoa(actualPort))),
	}
	s.mu.Unlock()

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.mu.Lock()
			s.server = nil
			s.status.Running = false
			s.status.LastError = err.Error()
			s.mu.Unlock()
		}
	}()
	return nil
}

func (s *TuneWebServer) Stop() error {
	s.mu.Lock()
	server := s.server
	if server == nil {
		s.status.Running = false
		s.mu.Unlock()
		return nil
	}
	s.server = nil
	s.status.Running = false
	s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *TuneWebServer) Status() TuneWebServerStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func preferredTuneWebAddress() string {
	items := telemetry.ListNetworkInterfaces()
	for _, item := range items {
		if item.IsUp && item.IsPrivate && !item.IsLoopback && item.Address != "0.0.0.0" && isPreferredWirelessName(item.Name+" "+item.DisplayName) {
			return item.Address
		}
	}
	for _, item := range items {
		if item.IsUp && item.IsPrivate && !item.IsLoopback && item.Address != "0.0.0.0" {
			return item.Address
		}
	}
	return "127.0.0.1"
}

func isPreferredWirelessName(value string) bool {
	name := strings.ToLower(value)
	return strings.Contains(name, "wlan") ||
		strings.Contains(name, "wi-fi") ||
		strings.Contains(name, "wifi") ||
		strings.Contains(name, "wireless") ||
		strings.Contains(name, "无线") ||
		strings.Contains(name, "無線")
}

func tuneWebRootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/tune", http.StatusFound)
}

func tuneWebPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(tuneWebHTML()))
}

func tuneWebGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var input storage.RoadStaticTuneBaselineInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeTuneWebJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON input"})
		return
	}
	result, err := storage.GenerateRoadStaticTuneBaseline(input)
	if err != nil {
		writeTuneWebJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeTuneWebJSON(w, http.StatusOK, result)
}

func writeTuneWebJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func tuneWebHTML() string {
	return `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
<title>FH6 远程快速调校</title>
<style>
:root{color-scheme:dark;--bg:#0e1316;--panel:#151b1f;--line:#2a343a;--text:#f6fbff;--muted:#9db2bd;--accent:#6ee7ff;--warn:#ffd166;--bad:#ff6b6b}
*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--text);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;font-size:16px;line-height:1.35}
main{width:min(760px,100%);margin:0 auto;padding:18px 14px 32px}.hero{padding:12px 0 18px}h1{font-size:24px;margin:0 0 6px}p{margin:0;color:var(--muted)}
.card{background:var(--panel);border:1px solid var(--line);border-radius:14px;padding:14px;margin:12px 0}.grid{display:grid;gap:12px}.row{display:grid;gap:8px}label span{display:block;color:var(--muted);font-size:13px;margin-bottom:6px}
input,select,button{width:100%;font:inherit;border-radius:12px;border:1px solid var(--line);background:#0d1215;color:var(--text);padding:13px 12px}button{background:linear-gradient(135deg,#0ea5c6,#2563eb);border:0;font-weight:800;margin-top:8px}button.secondary{background:#202a30}
input[type=range]{padding:0;height:34px;border:0;background:transparent;accent-color:var(--accent)}input[type=range]:disabled{opacity:.45}.slider-row{border:1px solid var(--line);border-radius:14px;padding:12px;background:#10171b}.slider-row.disabled{opacity:.65}.slider-head{display:flex;justify-content:space-between;gap:12px;align-items:center}.slider-head span{margin:0;color:var(--muted);font-size:13px}.slider-head strong{font-size:18px}.slider-scale{display:grid;grid-template-columns:1fr 1fr 1fr;gap:8px;color:var(--muted);font-size:12px}.slider-scale span:nth-child(2){text-align:center;color:var(--accent)}.slider-scale span:last-child{text-align:right}
.seg{display:grid;grid-template-columns:repeat(3,1fr);gap:8px}.seg button{background:#10171b;border:1px solid var(--line);margin:0}.seg button.active{background:#174151;border-color:var(--accent)}
.two{display:grid;grid-template-columns:1fr 1fr;gap:10px}.tire-size{display:grid;grid-template-columns:1fr auto 1fr auto 1fr;gap:8px;align-items:center}.tire-size b{color:var(--muted)}
.error{color:var(--bad);font-weight:700;margin-top:8px}.warn{color:var(--warn);font-weight:700;margin-top:8px}.hidden{display:none}
.result-group{margin-top:14px}.result-group h3{font-size:17px;margin:0 0 8px}.item{display:flex;justify-content:space-between;gap:12px;border-bottom:1px solid rgba(255,255,255,.08);padding:9px 0}.item span:first-child{color:var(--muted)}.item strong{text-align:right}
@media (min-width:700px){.grid{grid-template-columns:1fr 1fr}.wide{grid-column:1/-1}}
</style>
</head>
<body>
<main>
<section class="hero"><h1>FH6 远程快速调校</h1><p>仅生成预览，不复制结果，不保存档案。</p></section>
<section class="card">
<div class="grid">
<label class="row"><span>用途</span><select id="useCase"><option value="Road">公路</option><option value="Drift">漂移</option><option value="Rally">拉力</option><option value="Offroad">越野</option><option value="Drag">直线</option></select></label>
<label class="row"><span>轮胎类型</span><select id="tireCompound"><option value="sport">运动胎</option><option value="drift">漂移胎</option><option value="rally">拉力胎</option><option value="offroad">越野胎</option><option value="drag">直线胎</option><option value="stock">原厂/街胎</option><option value="semi">半热熔</option><option value="slick">光头胎</option></select></label>
<label class="row"><span>车重 KG</span><input id="weightKG" inputmode="numeric" placeholder="例如 1460"></label>
<label class="row"><span>前轮重量分配 %</span><input id="frontWeightPct" inputmode="numeric" placeholder="例如 60"></label>
<label class="row"><span>性能指数 PI</span><input id="pi" inputmode="numeric" placeholder="100-999"></label>
<div class="row"><span>驱动方式</span><div class="seg"><button type="button" data-drive="FWD">前驱</button><button type="button" data-drive="AWD">四驱</button><button type="button" data-drive="RWD" class="active">后驱</button></div></div>
<div class="row slider-row"><div class="slider-head"><span>平衡</span><strong id="balanceBiasValue">100</strong></div><input id="balanceBias" type="range" min="50" max="150" step="1" value="100"><div class="slider-scale"><span>稳定</span><span>中性</span><span>灵活</span></div></div>
<div class="row slider-row"><div class="slider-head"><span>硬度</span><strong id="stiffnessBiasValue">100</strong></div><input id="stiffnessBias" type="range" min="50" max="150" step="1" value="100"><div class="slider-scale"><span>软</span><span>中性</span><span>硬</span></div></div>
<div class="row slider-row disabled" id="speedBiasRow"><div class="slider-head"><span>速度</span><strong id="speedBiasValue">100</strong></div><input id="speedBias" type="range" min="50" max="150" step="1" value="100" disabled><div class="slider-scale"><span>极速</span><span>中性</span><span>加速</span></div></div>
</div>
<button class="secondary" type="button" id="toggleGearing">开启齿轮设置</button>
<div id="gearing" class="hidden">
<div class="grid">
<label class="row"><span>红线转速 RPM</span><input id="redlineRPM" inputmode="numeric"></label>
<label class="row"><span>档位数</span><input id="gearCount" inputmode="numeric" placeholder="2-10"></label>
<label class="row wide"><span>轮胎尺寸</span><div class="tire-size"><input id="tireWidth" inputmode="numeric" placeholder="245"><b>/</b><input id="tireAspect" inputmode="numeric" placeholder="35"><b>R</b><input id="tireRim" inputmode="numeric" placeholder="19"></div></label>
<label class="row wide"><span id="targetSpeedLabel">目标极速 km/h</span><input id="targetTopSpeedKmh" inputmode="numeric"></label>
</div>
</div>
<button type="button" id="generate">生成调校参数</button>
<div id="message"></div>
</section>
<section id="result" class="card hidden"></section>
</main>
<script>
const state={drivetrain:'RWD',gearing:false,hasPreview:false};
document.querySelectorAll('[data-drive]').forEach(btn=>btn.addEventListener('click',()=>{state.drivetrain=btn.dataset.drive;document.querySelectorAll('[data-drive]').forEach(b=>b.classList.toggle('active',b===btn));}));
function refreshUseCaseLabels(){const useCase=document.getElementById('useCase').value;document.getElementById('targetSpeedLabel').textContent=useCase==='Drift'?'目标漂移速度 km/h':useCase==='Drag'?'目标终点速度 km/h':'目标极速 km/h';}
document.getElementById('useCase').addEventListener('change',e=>{const defaults={Road:'sport',Drift:'drift',Rally:'rally',Offroad:'offroad',Drag:'drag'};const tire=document.getElementById('tireCompound');if(['sport','drift','rally','offroad','drag'].includes(tire.value))tire.value=defaults[e.target.value]||'sport';refreshUseCaseLabels();});
refreshUseCaseLabels();
function snapBiasValue(value){value=Math.max(50,Math.min(150,Math.round(Number(value)||100)));return Math.abs(value-100)<=3?100:value}
function setBiasValue(id,value){const input=document.getElementById(id);const snapped=snapBiasValue(value);input.value=String(snapped);document.getElementById(id+'Value').textContent=String(snapped)}
function refreshSpeedBiasEnabled(){const input=document.getElementById('speedBias');const row=document.getElementById('speedBiasRow');input.disabled=!state.gearing;row.classList.toggle('disabled',!state.gearing);if(!state.gearing)setBiasValue('speedBias',100)}
async function regenerateIfPreview(){if(!state.hasPreview)return;await generatePreview()}
['balanceBias','stiffnessBias','speedBias'].forEach(id=>{const input=document.getElementById(id);input.addEventListener('input',()=>{document.getElementById(id+'Value').textContent=String(input.value)});input.addEventListener('change',()=>{setBiasValue(id,input.value);regenerateIfPreview()});});
document.getElementById('toggleGearing').addEventListener('click',()=>{state.gearing=!state.gearing;document.getElementById('gearing').classList.toggle('hidden',!state.gearing);document.getElementById('toggleGearing').textContent=state.gearing?'关闭齿轮设置':'开启齿轮设置';refreshSpeedBiasEnabled();if(!state.gearing)regenerateIfPreview();});
refreshSpeedBiasEnabled();
function intValue(id){const raw=document.getElementById(id).value.trim();if(!/^\d+$/.test(raw))throw new Error(id+' 必须填写整数');return Number(raw)}
function tireDiameterCm(){const w=intValue('tireWidth'),a=intValue('tireAspect'),r=intValue('tireRim');return (r*25.4+2*w*(a/100))/10}
function buildInput(){const input={useCase:document.getElementById('useCase').value,tireCompound:document.getElementById('tireCompound').value,drivetrain:state.drivetrain,weightKG:intValue('weightKG'),frontWeightPct:intValue('frontWeightPct'),pi:intValue('pi'),balanceBias:intValue('balanceBias'),stiffnessBias:intValue('stiffnessBias'),speedBias:intValue('speedBias'),frontRideHeightAdjustable:true,rearRideHeightAdjustable:true,frontAeroAdjustable:true,rearAeroAdjustable:true};if(state.gearing){input.redlineRPM=intValue('redlineRPM');input.gearCount=intValue('gearCount');input.tireDiameterCm=tireDiameterCm();input.targetTopSpeedKmh=intValue('targetTopSpeedKmh')}return input}
function label(key){const map={frontTirePressure:'前侧胎压',rearTirePressure:'后侧胎压',finalDrive:'终传比',gear1:'1 挡',gear2:'2 挡',gear3:'3 挡',gear4:'4 挡',gear5:'5 挡',gear6:'6 挡',gear7:'7 挡',gear8:'8 挡',gear9:'9 挡',gear10:'10 挡',frontCamber:'前侧外倾角',rearCamber:'后侧外倾角',frontToe:'前侧束角',rearToe:'后侧束角',caster:'前轮后倾角',frontArb:'前侧防倾杆',rearArb:'后侧防倾杆',frontSpring:'前侧弹簧',rearSpring:'后侧弹簧',frontRideHeight:'前侧车身高度',rearRideHeight:'后侧车身高度',frontRebound:'前侧回弹硬度',rearRebound:'后侧回弹硬度',frontBump:'前侧压缩硬度',rearBump:'后侧压缩硬度',frontAero:'前侧下压力',rearAero:'后侧下压力',brakeBalance:'制动力平衡',brakePressure:'制动力压力',frontDiffAccel:'前侧加速',frontDiffDecel:'前侧减速',rearDiffAccel:'后侧加速',rearDiffDecel:'后侧减速',centerDiffBalance:'中央平衡'};return map[key]||key}
function groupName(group){return {tire:'轮胎',gearing:'齿比',alignment:'轮胎定位',antiroll:'防倾杆',springs:'弹簧与车高',damping:'阻尼',aero:'空气动力学设置',brake:'刹车',differential:'差速器',power:'车辆信息'}[group]||group}
const fieldOrder=['frontTirePressure','rearTirePressure','finalDrive','gear1','gear2','gear3','gear4','gear5','gear6','gear7','gear8','gear9','gear10','frontCamber','rearCamber','frontToe','rearToe','caster','frontArb','rearArb','frontSpring','rearSpring','frontRideHeight','rearRideHeight','frontRebound','rearRebound','frontBump','rearBump','brakeBalance','brakePressure','frontAero','rearAero','frontDiffAccel','frontDiffDecel','rearDiffAccel','rearDiffDecel','centerDiffBalance'];
function fieldRank(key){const i=fieldOrder.indexOf(key);return i>=0?i:9999}
function render(data){const box=document.getElementById('result');const groups={};(data.generatedFields||[]).forEach(f=>{(groups[f.group]||(groups[f.group]=[])).push(f)});(data.tierRecommendations||[]).forEach(f=>{(groups[f.group]||(groups[f.group]=[])).push({...f,value:f.tier,unit:''})});Object.keys(groups).forEach(g=>groups[g].sort((a,b)=>fieldRank(a.fieldKey)-fieldRank(b.fieldKey)));const order=['tire','gearing','alignment','antiroll','springs','damping','brake','aero','differential','power'];box.innerHTML='<h2>调校参数预览</h2>'+order.filter(g=>groups[g]).map(g=>'<div class="result-group"><h3>'+groupName(g)+'</h3>'+groups[g].map(f=>'<div class="item"><span>'+label(f.fieldKey)+'</span><strong>'+((f.value??'--')+(f.unit?' '+f.unit:''))+'</strong></div>').join('')+'</div>').join('');box.classList.remove('hidden');state.hasPreview=true}
async function generatePreview(){const msg=document.getElementById('message');msg.className='';msg.textContent='';try{const res=await fetch('/api/tune/generate',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(buildInput())});const data=await res.json();if(!res.ok)throw new Error(data.error||'生成失败');render(data)}catch(err){msg.className='error';msg.textContent=err.message||String(err)}}
document.getElementById('generate').addEventListener('click',generatePreview);
</script>
</body>
</html>`
}
